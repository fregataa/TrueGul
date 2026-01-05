import asyncio
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.services.detector import detector_service
from app.services.feedback import feedback_service
from app.services.callback import callback_client
from app.services.consumer import create_task_processor

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    detector_service.load_model()
    feedback_service.init_client()

    task_processor = create_task_processor()
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
async def health_check():
    return {"status": "healthy"}


@app.get("/")
async def root():
    return {"message": "TrueGul ML Server"}
