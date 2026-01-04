import logging
from openai import OpenAI

from app.config import settings
from app.schemas.task import WritingType

logger = logging.getLogger(__name__)

PROMPTS = {
    WritingType.ESSAY: """You are a writing coach providing feedback on an essay.
Analyze the following text and provide constructive feedback focusing on:
- Structure and organization
- Clarity of arguments
- Writing style and tone
- Grammar and vocabulary

Be concise and actionable. Provide feedback in Korean.

Text:
{text}

AI Detection Score: {ai_score}% (higher means more likely AI-generated)

Provide your feedback:""",
    WritingType.COVER_LETTER: """You are a career coach providing feedback on a cover letter.
Analyze the following text and provide constructive feedback focusing on:
- Professional tone and presentation
- Relevance to job application context
- Authenticity and personal voice
- Structure and persuasiveness

Be concise and actionable. Provide feedback in Korean.

Text:
{text}

AI Detection Score: {ai_score}% (higher means more likely AI-generated)

Provide your feedback:""",
}


class FeedbackService:
    def __init__(self):
        self._client = None

    def init_client(self):
        if settings.openai_api_key:
            self._client = OpenAI(api_key=settings.openai_api_key)
            logger.info("OpenAI client initialized")
        else:
            logger.warning("OpenAI API key not set, feedback will be unavailable")

    def generate_feedback(self, text: str, writing_type: WritingType, ai_score: float) -> str:
        if self._client is None:
            return "피드백 서비스를 사용할 수 없습니다. (OpenAI API 키 미설정)"

        prompt = PROMPTS.get(writing_type, PROMPTS[WritingType.ESSAY])
        formatted_prompt = prompt.format(text=text[:2000], ai_score=f"{ai_score:.1f}")

        response = self._client.chat.completions.create(
            model="gpt-4o-mini",
            messages=[{"role": "user", "content": formatted_prompt}],
            max_tokens=500,
            temperature=0.7,
        )

        return response.choices[0].message.content or "피드백을 생성할 수 없습니다."


feedback_service = FeedbackService()
