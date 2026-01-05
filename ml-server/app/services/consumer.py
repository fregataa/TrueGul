import logging
import time

from app.config import settings
from app.mq import Consumer, Message, RedisConsumer
from app.schemas.task import (
    AnalysisTask,
    AnalysisCallback,
    AnalysisResult,
    AnalysisError,
    ErrorCode,
)
from app.services.detector import AIDetectorService
from app.services.feedback import FeedbackService
from app.services.callback import CallbackClient

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

        task: AnalysisTask | None = None
        try:
            task = AnalysisTask.model_validate_json(task_data)
            logger.info(f"Processing task: {task.task_id}")

            start_time = time.time()

            ai_score = self._detector.detect(task.content)

            feedback = self._feedback.generate_feedback(task.content, task.writing_type, ai_score)

            latency_ms = int((time.time() - start_time) * 1000)

            callback = AnalysisCallback(
                task_id=task.task_id,
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
                task_id=task.task_id if task else "unknown",
                status="failed",
                error=AnalysisError(
                    code=error_code,
                    message=str(e),
                    retryable=retryable,
                ),
            )

        if task:
            await self._callback.send_callback(task.callback_url, callback)
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
