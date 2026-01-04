import logging
from transformers import pipeline

from app.config import settings

logger = logging.getLogger(__name__)


class AIDetectorService:
    def __init__(self):
        self._pipeline = None

    def load_model(self):
        logger.info(f"Loading model: {settings.model_name}")
        self._pipeline = pipeline(
            "text-classification",
            model=settings.model_name,
            device=-1,
        )
        logger.info("Model loaded successfully")

    def detect(self, text: str) -> float:
        if self._pipeline is None:
            raise RuntimeError("Model not loaded. Call load_model() first.")

        result = self._pipeline(text, truncation=True, max_length=512)
        label = result[0]["label"]
        score = result[0]["score"]

        if label.lower() in ("chatgpt", "ai", "fake", "machine"):
            return score * 100
        return (1 - score) * 100


detector_service = AIDetectorService()
