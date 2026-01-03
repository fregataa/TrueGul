from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

router = APIRouter()


class AnalyzeRequest(BaseModel):
    text: str
    writing_type: str = "essay"


class AnalyzeResponse(BaseModel):
    ai_score: float
    feedback: str
    model_version: str


@router.post("/analyze", response_model=AnalyzeResponse)
async def analyze_text(request: AnalyzeRequest):
    if not request.text.strip():
        raise HTTPException(status_code=400, detail="Text cannot be empty")

    # Placeholder response - to be implemented
    return AnalyzeResponse(
        ai_score=0.0,
        feedback="Analysis pending implementation",
        model_version="0.1.0",
    )
