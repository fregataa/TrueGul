from dataclasses import dataclass

from pydantic_settings import BaseSettings

from app.common.types import ModelLoaderType


@dataclass(frozen=True)
class ModelLoaderConfig:
    """Configuration for model loader."""

    loader_type: ModelLoaderType
    bucket: str
    local_dir: str


@dataclass(frozen=True)
class DetectorModelConfig:
    """Configuration for AI detector model."""

    name: str  # HuggingFace model name (fallback)
    s3_key: str
    version: str


@dataclass(frozen=True)
class FeedbackModelConfig:
    """Configuration for feedback generation model."""

    s3_key: str
    n_ctx: int
    n_threads: int
    max_tokens: int


class Settings(BaseSettings):
    port: int = 8000
    redis_url: str = "redis://localhost:6379"
    api_server_url: str = "http://localhost:8080"
    ml_callback_secret: str = ""

    # Model Loader Configuration
    ml_model_loader_type: ModelLoaderType = ModelLoaderType.S3
    ml_models_bucket: str = ""  # Empty = local/HuggingFace mode
    ml_models_local_dir: str = "/app/models"

    # AI Detector Model
    detector_model_name: str = "Hello-SimpleAI/chatgpt-detector-roberta"
    detector_model_s3_key: str = "detector/chatgpt-detector-roberta"
    detector_model_version: str = "0.1.0"

    # Feedback Model
    feedback_model_s3_key: str = "feedback/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf"
    feedback_n_ctx: int = 2048
    feedback_n_threads: int = 4
    feedback_max_tokens: int = 256

    # MQ Configuration
    stream_name: str = "analysis_tasks"
    consumer_group: str = "ml_workers"
    consumer_name: str = "worker_1"
    max_retries: int = 3

    class Config:
        env_file = ".env"

    def get_model_loader_config(self) -> ModelLoaderConfig:
        return ModelLoaderConfig(
            loader_type=self.ml_model_loader_type,
            bucket=self.ml_models_bucket,
            local_dir=self.ml_models_local_dir,
        )

    def get_detector_model_config(self) -> DetectorModelConfig:
        return DetectorModelConfig(
            name=self.detector_model_name,
            s3_key=self.detector_model_s3_key,
            version=self.detector_model_version,
        )

    def get_feedback_model_config(self) -> FeedbackModelConfig:
        return FeedbackModelConfig(
            s3_key=self.feedback_model_s3_key,
            n_ctx=self.feedback_n_ctx,
            n_threads=self.feedback_n_threads,
            max_tokens=self.feedback_max_tokens,
        )


settings = Settings()
