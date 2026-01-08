import asyncio
import contextlib
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware

from app.config import settings
from app.model_adapter import ModelAdapter
from app.model_loader.creator import create_model_loaders
from app.services.callback import CallbackClient
from app.services.consumer import create_task_processor

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Initialize model loaders
    loaders = create_model_loaders(settings.get_model_loader_config())

    # Initialize model adapter
    adapter = ModelAdapter(loaders)

    # Create services with loaded models
    app.state.detector = adapter.create_detector_service(settings.get_detector_model_config())
    app.state.feedback = adapter.create_feedback_service(settings.get_feedback_model_config())

    callback_client = CallbackClient()

    task_processor = create_task_processor(app.state.detector, app.state.feedback, callback_client)
    await task_processor.connect()

    consumer_task = asyncio.create_task(task_processor.start())

    yield

    await task_processor.stop()
    consumer_task.cancel()
    with contextlib.suppress(asyncio.CancelledError):
        await consumer_task
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

    if detector is not None:
        services["detector_model"] = "healthy"
    else:
        services["detector_model"] = "not loaded"
        overall_status = "unhealthy"

    if feedback is not None:
        services["feedback_model"] = "healthy"
    else:
        services["feedback_model"] = "not loaded"
        overall_status = "unhealthy"

    return {"status": overall_status, "services": services}


@app.get("/health/live")
async def liveness():
    return {"status": "alive"}


@app.get("/health/ready")
async def readiness(request: Request):
    return await health_check(request)


@app.get("/")
async def root():
    return {"message": "TrueGul ML Server"}
