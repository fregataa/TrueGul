from enum import StrEnum
from typing import Optional

from pydantic import BaseModel


class WritingType(StrEnum):
    ESSAY = "essay"
    COVER_LETTER = "cover_letter"


class ErrorCode(StrEnum):
    ML_MODEL_ERROR = "ML_MODEL_ERROR"
    OPENAI_API_ERROR = "OPENAI_API_ERROR"
    INVALID_INPUT = "INVALID_INPUT"
    TIMEOUT = "TIMEOUT"
    INTERNAL_ERROR = "INTERNAL_ERROR"


class AnalysisTask(BaseModel):
    version: str = "1"
    task_id: str
    writing_id: str
    content: str
    writing_type: WritingType
    callback_url: str


class AnalysisResult(BaseModel):
    ai_probability: float
    feedback: str
    latency_ms: int


class AnalysisError(BaseModel):
    code: ErrorCode
    message: str
    retryable: bool


class AnalysisCallback(BaseModel):
    version: str = "1"
    task_id: str
    status: str
    result: Optional[AnalysisResult] = None
    error: Optional[AnalysisError] = None
