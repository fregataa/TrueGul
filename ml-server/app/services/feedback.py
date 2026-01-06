import logging
import textwrap

from llama_cpp import Llama

from app.config import settings
from app.schemas.task import WritingType

logger = logging.getLogger(__name__)


class FeedbackService:
    def __init__(self):
        self._llm: Llama | None = None

    def load_model(self):
        logger.info(f"Loading feedback model: {settings.feedback_model_path}")
        self._llm = Llama(
            model_path=settings.feedback_model_path,
            n_ctx=settings.feedback_n_ctx,
            n_threads=settings.feedback_n_threads,
            verbose=False,
        )
        logger.info("Feedback model loaded successfully")

    def is_loaded(self) -> bool:
        return self._llm is not None

    def generate_feedback(
        self, text: str, writing_type: WritingType, ai_score: float
    ) -> str:
        if self._llm is None:
            raise RuntimeError("Model not loaded")

        prompt = self._build_prompt(text, writing_type, ai_score)

        output = self._llm(
            prompt,
            max_tokens=settings.feedback_max_tokens,
            temperature=0.7,
            stop=["</s>", "\n\n\n"],
        )

        return output["choices"][0]["text"].strip()

    def _detect_language(self, text: str) -> str:
        """Simple heuristic: if > 30% Korean characters, it's Korean"""
        korean_count = sum(1 for c in text if "\uac00" <= c <= "\ud7a3")
        return "ko" if korean_count / max(len(text), 1) > 0.3 else "en"

    def _build_prompt(
        self, text: str, writing_type: WritingType, ai_score: float
    ) -> str:
        lang = self._detect_language(text)

        if lang == "ko":
            writing_type_str = (
                "에세이" if writing_type == WritingType.ESSAY else "자기소개서"
            )
            instruction = (
                f"다음 {writing_type_str}에 대해 2-3문장으로 피드백을 제공해주세요."
            )
            system_msg = (
                "You are a writing coach providing brief, constructive feedback in Korean."
            )
        else:
            writing_type_str = (
                "essay" if writing_type == WritingType.ESSAY else "cover letter"
            )
            instruction = (
                f"Provide 2-3 sentences of feedback for this {writing_type_str}."
            )
            system_msg = (
                "You are a writing coach providing brief, constructive feedback in English."
            )

        return textwrap.dedent(f"""\
            <|system|>
            {system_msg}
            </s>
            <|user|>
            {instruction}
            AI Score: {ai_score:.1f}%

            Text:
            {text[:1500]}
            </s>
            <|assistant|>
            """)
