import logging

from llama_cpp import Llama
from transformers import pipeline
from transformers.pipelines import Pipeline

from app.config import DetectorModelConfig, FeedbackModelConfig
from app.model_loader.creator import ModelLoaders
from app.services.detector import AIDetectorService
from app.services.feedback import FeedbackService

logger = logging.getLogger(__name__)


class ModelAdapter:
    """Adapter for loading models and creating services."""

    def __init__(self, loaders: ModelLoaders) -> None:
        self._loader = loaders.primary
        self._fallback_loader = loaders.fallback

    def _ensure_file(self, key: str) -> str:
        """Try to ensure file from primary loader, then fallback."""
        model_path = self._loader.ensure_file(key)
        if model_path:
            return model_path

        if self._fallback_loader:
            logger.info(f"Primary loader failed, trying fallback for file: {key}")
            model_path = self._fallback_loader.ensure_file(key)

        return model_path

    def _ensure_directory(self, key: str) -> str:
        """Try to ensure directory from primary loader, then fallback."""
        model_path = self._loader.ensure_directory(key)
        if model_path:
            return model_path

        if self._fallback_loader:
            logger.info(f"Primary loader failed, trying fallback for directory: {key}")
            model_path = self._fallback_loader.ensure_directory(key)

        return model_path

    def _load_detector_pipeline(self, config: DetectorModelConfig) -> Pipeline:
        """Load AI detector pipeline from loader or HuggingFace hub."""
        model_path = self._ensure_directory(config.s3_key)

        if model_path:
            logger.info(f"Using local detector model: {model_path}")
        else:
            model_path = config.name
            logger.info(f"Loading detector model from HuggingFace hub: {model_path}")

        return pipeline("text-classification", model=model_path, device=-1)

    def _load_feedback_llm(self, config: FeedbackModelConfig) -> Llama:
        """Load feedback Llama model from loader."""
        model_path = self._ensure_file(config.s3_key)

        if not model_path:
            raise RuntimeError(
                f"Feedback model not found. "
                f"Set ML_MODELS_BUCKET or place model at local_dir/{config.s3_key}"
            )

        logger.info(f"Loading feedback model: {model_path}")
        return Llama(
            model_path=model_path,
            n_ctx=config.n_ctx,
            n_threads=config.n_threads,
            verbose=False,
        )

    def create_detector_service(self, config: DetectorModelConfig) -> AIDetectorService:
        """Create AIDetectorService with loaded model."""
        logger.info("Loading AI detector model...")
        detector_pipeline = self._load_detector_pipeline(config)
        logger.info("AI detector model loaded successfully")
        return AIDetectorService(pipeline=detector_pipeline)

    def create_feedback_service(self, config: FeedbackModelConfig) -> FeedbackService:
        """Create FeedbackService with loaded model."""
        logger.info("Loading feedback model...")
        feedback_llm = self._load_feedback_llm(config)
        logger.info("Feedback model loaded successfully")
        return FeedbackService(llm=feedback_llm, max_tokens=config.max_tokens)
