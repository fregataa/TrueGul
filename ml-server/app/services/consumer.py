import logging
import time

from app.config import settings
from app.mq import Consumer, Message, RedisConsumer
from app.schemas.task import (
    AnalysisCallback,
    AnalysisError,
    AnalysisResult,
    AnalysisTask,
    ErrorCode,
)
from app.services.callback import CallbackClient
from app.services.detector import AIDetectorService
from app.services.feedback import FeedbackService

logger = logging.getLogger(__name__)


class TaskProcessor:
    def __init__(
        self,
        consumer: Consumer,
        detector: AIDetectorService,
        feedback: FeedbackService,
        callback: CallbackClient,
    ) -> None:
        self._consumer = consumer
        self._detector = detector
        self._feedback = feedback
        self._callback = callback

    async def connect(self) -> None:
        await self._consumer.connect()

    async def disconnect(self) -> None:
        await self._consumer.disconnect()

    async def start(self) -> None:
        await self._consumer.start(self._handle_message)

    async def stop(self) -> None:
        await self._consumer.stop()

    async def _handle_message(self, message: Message) -> None:
        task_data = message.data.get("task")
        if not task_data:
            logger.error(f"Invalid message format: {message.data}")
            await self._consumer.ack(message.id)
            return

        task_id: str = "unknown"
        callback_url: str | None = None
        callback: AnalysisCallback

        try:
            task = AnalysisTask.model_validate_json(task_data)
            task_id = task.task_id
            callback_url = task.callback_url
            content = task.content
            writing_type = task.writing_type
            logger.info(f"Processing task: {task_id}")

            start_time = time.time()

            ai_score = self._detector.detect(content)

            feedback = self._feedback.generate_feedback(content, writing_type, ai_score)

            latency_ms = int((time.time() - start_time) * 1000)

            callback = AnalysisCallback(
                task_id=task_id,
                status="completed",
                result=AnalysisResult(
                    ai_probability=ai_score,
                    feedback=feedback,
                    latency_ms=latency_ms,
                ),
            )

        except Exception as e:
            logger.exception(f"Error processing task: {e}")
            error_code = ErrorCode.INTERNAL_ERROR
            retryable = False

            if "model" in str(e).lower():
                error_code = ErrorCode.ML_MODEL_ERROR
                retryable = True
            elif "openai" in str(e).lower():
                error_code = ErrorCode.OPENAI_API_ERROR
                retryable = True

            callback = AnalysisCallback(
                task_id=task_id,
                status="failed",
                error=AnalysisError(
                    code=error_code,
                    message=str(e),
                    retryable=retryable,
                ),
            )

        if callback_url:
            await self._callback.send_callback(callback_url, callback)
        await self._consumer.ack(message.id)


def create_task_processor(
    detector: AIDetectorService,
    feedback: FeedbackService,
    callback: CallbackClient,
) -> TaskProcessor:
    consumer = RedisConsumer(
        redis_url=settings.redis_url,
        stream_name=settings.stream_name,
        consumer_group=settings.consumer_group,
        consumer_name=settings.consumer_name,
    )
    return TaskProcessor(consumer, detector, feedback, callback)
