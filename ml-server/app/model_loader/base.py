from __future__ import annotations

from abc import ABC, abstractmethod
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from app.config import ModelLoaderConfig


class ModelLoader(ABC):
    """Abstract base class for loading ML models from various sources."""

    @classmethod
    @abstractmethod
    def create(cls, config: ModelLoaderConfig) -> ModelLoader:
        """Factory method to create a ModelLoader instance from configuration."""
        pass

    @abstractmethod
    def is_enabled(self) -> bool:
        """Check if this loader is enabled and available."""
        pass

    @abstractmethod
    def ensure_file(self, key: str) -> str:
        """
        Ensure a single model file is available locally.

        Args:
            key: Identifier for the model file (e.g., S3 key or path)

        Returns:
            Local path to the model file, or empty string if not available
        """
        pass

    @abstractmethod
    def ensure_directory(self, key: str) -> str:
        """
        Ensure a model directory is available locally.

        Args:
            key: Identifier for the model directory (e.g., S3 prefix)

        Returns:
            Local path to the model directory, or empty string if not available
        """
        pass
