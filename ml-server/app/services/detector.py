import logging

from transformers.pipelines import Pipeline

logger = logging.getLogger(__name__)


class AIDetectorService:
    def __init__(self, pipeline: Pipeline):
        self._pipeline = pipeline

    def detect(self, text: str) -> float:
        result = self._pipeline(text, truncation=True, max_length=512)
        label = result[0]["label"]
        score = result[0]["score"]

        if label.lower() in ("chatgpt", "ai", "fake", "machine"):
            return score * 100
        return (1 - score) * 100
