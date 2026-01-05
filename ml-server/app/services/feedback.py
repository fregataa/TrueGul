import logging

from app.schemas.task import WritingType

logger = logging.getLogger(__name__)


class FeedbackService:
    """Placeholder for feedback service. Will be implemented with local model in v1."""

    def __init__(self):
        pass

    def init_client(self):
        logger.info("Feedback service initialized (placeholder - no model loaded)")

    def is_initialized(self) -> bool:
        return True

    def generate_feedback(self, text: str, writing_type: WritingType, ai_score: float) -> str:
        # TODO: Implement with local model in v1
        return ""
