from enum import StrEnum


class ModelLoaderType(StrEnum):
    """Supported model loader types."""

    S3 = "s3"
    LOCAL = "local"
