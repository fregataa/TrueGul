from __future__ import annotations

import logging
from pathlib import Path
from typing import TYPE_CHECKING, override

from app.model_loader.base import ModelLoader

if TYPE_CHECKING:
    from app.config import ModelLoaderConfig

logger = logging.getLogger(__name__)


class LocalModelLoader(ModelLoader):
    """ModelLoader implementation that loads models from local filesystem."""

    @classmethod
    @override
    def create(cls, config: ModelLoaderConfig) -> LocalModelLoader:
        return cls(local_dir=config.local_dir)

    def __init__(self, local_dir: str = "/app/models"):
        self._local_dir = Path(local_dir)
        logger.info(f"LocalModelLoader initialized with local_dir: {self._local_dir}")

    @override
    def is_enabled(self) -> bool:
        return True

    @override
    def ensure_file(self, key: str) -> str:
        local_path = self._local_dir / key

        if local_path.exists() and local_path.is_file():
            logger.info(f"Model file found at {local_path}")
            return str(local_path)

        logger.info(f"Model file not found: {local_path}")
        return ""

    @override
    def ensure_directory(self, key: str) -> str:
        local_path = self._local_dir / key

        if local_path.exists() and local_path.is_dir():
            logger.info(f"Model directory found at {local_path}")
            return str(local_path)

        logger.info(f"Model directory not found: {local_path}")
        return ""
