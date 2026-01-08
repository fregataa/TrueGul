from __future__ import annotations

import logging
from pathlib import Path
from typing import TYPE_CHECKING, override

import boto3

from app.model_loader.base import ModelLoader

if TYPE_CHECKING:
    from app.config import ModelLoaderConfig

logger = logging.getLogger(__name__)


class S3ModelLoader(ModelLoader):
    """ModelLoader implementation that downloads models from S3."""

    @classmethod
    @override
    def create(cls, config: ModelLoaderConfig) -> S3ModelLoader:
        return cls(bucket=config.bucket, local_dir=config.local_dir)

    def __init__(self, bucket: str = "", local_dir: str = "/app/models"):
        self._bucket = bucket
        self._local_dir = Path(local_dir)
        self._s3_client = None

        if self._bucket:
            logger.info(f"S3ModelLoader initialized with bucket: {self._bucket}")
        else:
            logger.info("S3ModelLoader initialized in local-only mode (no S3)")

    @property
    def _s3(self):
        """Lazy initialization of S3 client."""
        if self._s3_client is None and self._bucket:

            self._s3_client = boto3.client("s3")
        return self._s3_client

    @override
    def is_enabled(self) -> bool:
        return bool(self._bucket)

    @override
    def ensure_file(self, key: str) -> str:
        local_path = self._local_dir / key

        if local_path.exists():
            logger.info(f"Model file already exists at {local_path}")
            return str(local_path)

        if not self.is_enabled():
            logger.info(f"S3 not enabled, file not found locally: {local_path}")
            return ""

        self._download_file(key, local_path)
        return str(local_path)

    @override
    def ensure_directory(self, key: str) -> str:
        local_path = self._local_dir / key

        if local_path.exists():
            logger.info(f"Model directory already exists at {local_path}")
            return str(local_path)

        if not self.is_enabled():
            logger.info(f"S3 not enabled, directory not found locally: {local_path}")
            return ""

        self._download_directory(key, local_path)
        return str(local_path)

    def _download_file(self, s3_key: str, local_path: Path) -> None:
        """Download a single file from S3."""
        logger.info(f"Downloading s3://{self._bucket}/{s3_key} to {local_path}")
        local_path.parent.mkdir(parents=True, exist_ok=True)
        self._s3.download_file(self._bucket, s3_key, str(local_path))
        logger.info(f"Downloaded {s3_key} successfully")

    def _download_directory(self, s3_prefix: str, local_path: Path) -> None:
        """Download all files under an S3 prefix (directory)."""
        logger.info(f"Downloading s3://{self._bucket}/{s3_prefix}/ to {local_path}")
        local_path.mkdir(parents=True, exist_ok=True)

        paginator = self._s3.get_paginator("list_objects_v2")
        for page in paginator.paginate(Bucket=self._bucket, Prefix=s3_prefix):
            for obj in page.get("Contents", []):
                s3_key = obj["Key"]
                if s3_key.endswith("/"):
                    continue

                relative_path = s3_key[len(s3_prefix) :].lstrip("/")
                if not relative_path:
                    continue

                file_local_path = local_path / relative_path
                file_local_path.parent.mkdir(parents=True, exist_ok=True)

                logger.debug(f"Downloading {s3_key} -> {file_local_path}")
                self._s3.download_file(self._bucket, s3_key, str(file_local_path))

        logger.info(f"Downloaded {s3_prefix}/ successfully")
