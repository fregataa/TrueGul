from dataclasses import dataclass

from app.common.types import ModelLoaderType
from app.config import ModelLoaderConfig

from .base import ModelLoader
from .local import LocalModelLoader
from .s3 import S3ModelLoader


@dataclass(frozen=True)
class ModelLoaders:
    """Container for primary and fallback model loaders."""

    primary: ModelLoader
    fallback: ModelLoader | None = None


def create_model_loaders(config: ModelLoaderConfig) -> ModelLoaders:
    """Factory function to create ModelLoaders based on configuration."""
    fallback = LocalModelLoader(local_dir=config.local_dir)

    match config.loader_type:
        case ModelLoaderType.S3:
            primary = S3ModelLoader.create(config)
        case ModelLoaderType.LOCAL:
            primary = LocalModelLoader.create(config)
            fallback = None  # No fallback needed for local loader
        case _:
            raise ValueError(f"Unsupported model loader type: {config.loader_type}")

    return ModelLoaders(primary=primary, fallback=fallback)
