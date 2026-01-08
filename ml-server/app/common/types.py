
from enum import Enum, StrEnum


class ModelLoaderType(StrEnum):
    """Supported model loader types."""

    S3 = "s3"
    LOCAL = "local"
