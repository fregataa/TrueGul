import asyncio
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware

from app.services.detector import AIDetectorService
from app.services.feedback import FeedbackService
from app.services.callback import CallbackClient
from app.services.consumer import create_task_processor

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    app.state.detector = AIDetectorService()
    app.state.detector.load_model()

    app.state.feedback = FeedbackService()
    app.state.feedback.init_client()

    callback_client = CallbackClient()

    task_processor = create_task_processor(
        app.state.detector, app.state.feedback, callback_client
    )
    await task_processor.connect()

    consumer_task = asyncio.create_task(task_processor.start())

    yield

    await task_processor.stop()
    consumer_task.cancel()
    try:
        await consumer_task
    except asyncio.CancelledError:
        pass
    await task_processor.disconnect()
    await callback_client.close()


app = FastAPI(
    title="TrueGul ML Server",
    description="AI writing detection and feedback service",
    version="0.1.0",
    lifespan=lifespan,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/health")
async def health_check(request: Request):
    services = {}
    overall_status = "healthy"

    detector = getattr(request.app.state, "detector", None)
    feedback = getattr(request.app.state, "feedback", None)

    # Check model status
    if detector is not None and detector.is_loaded():
        services["model"] = "healthy"
    else:
        services["model"] = "not loaded"
        overall_status = "unhealthy"

    # Check OpenAI client status
    if feedback is not None and feedback.is_initialized():
        services["openai"] = "healthy"
    else:
        services["openai"] = "not initialized"
        overall_status = "unhealthy"

    return {
        "status": overall_status,
        "services": services
    }


@app.get("/health/live")
async def liveness():
    return {"status": "alive"}


@app.get("/health/ready")
async def readiness(request: Request):
    return await health_check(request)


@app.get("/")
async def root():
    return {"message": "TrueGul ML Server"}
