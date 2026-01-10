# Sprint 2: Backend - ML Server

## Overview

| 항목 | 내용 |
|------|------|
| 목표 | LLM 기반 TOPIK 채점 구현 |
| 선행 조건 | Sprint 0 완료 |
| 병렬 진행 | Sprint 1 (API Server) |
| 후속 Sprint | Sprint 4 (Integration) |

---

## Tasks

| ID | Task | 파일 | 상태 | 비고 |
|----|------|------|------|------|
| S2-1 | LLM Client 추상화 | `app/clients/llm_client.py` | TODO | |
| S2-2 | Claude Client 구현 | `app/clients/claude_client.py` | TODO | |
| S2-3 | OpenAI Client 구현 | `app/clients/openai_client.py` | TODO | |
| S2-4 | 프롬프트 빌더 | `app/prompts/topik_scoring.py` | TODO | |
| S2-5 | Few-shot 데이터 로더 | `app/data/few_shot_loader.py` | TODO | |
| S2-6 | TOPIK 채점 서비스 | `app/services/topik_scorer.py` | TODO | |
| S2-7 | Worker 수정 | `app/worker.py` | TODO | |
| S2-8 | AI 감지 결과 통합 | `app/services/detector.py` | TODO | 기존 코드 활용 |
| S2-9 | Callback 페이로드 수정 | `app/services/callback.py` | TODO | |
| S2-10 | Config 업데이트 | `app/config.py` | TODO | |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        ML Server                            │
│                                                             │
│  ┌───────────┐    ┌───────────┐    ┌───────────────────┐   │
│  │   Redis   │───▶│  Worker   │───▶│ Callback Sender   │   │
│  │ Consumer  │    │           │    │                   │   │
│  └───────────┘    └─────┬─────┘    └───────────────────┘   │
│                         │                                   │
│         ┌───────────────┼───────────────┐                  │
│         ▼               ▼               ▼                  │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐          │
│  │ AI Detector │ │TOPIK Scorer │ │   Prompt    │          │
│  │  (RoBERTa)  │ │             │ │   Builder   │          │
│  └─────────────┘ └──────┬──────┘ └─────────────┘          │
│                         │                                   │
│                         ▼                                   │
│                  ┌─────────────┐                           │
│                  │ LLM Client  │                           │
│                  │ (Claude/GPT)│                           │
│                  └─────────────┘                           │
└─────────────────────────────────────────────────────────────┘
```

---

## S2-1: LLM Client Abstraction

```python
# app/clients/llm_client.py
from abc import ABC, abstractmethod
from typing import TypedDict

class ScoringResult(TypedDict):
    scores: dict[str, int]
    feedback: dict[str, str]
    level_estimate: str

class LLMClient(ABC):
    @abstractmethod
    async def score_essay(
        self,
        content: str,
        prompt: str,
        few_shot_examples: list[dict]
    ) -> ScoringResult:
        """Score a TOPIK essay and return structured result."""
        pass

    @abstractmethod
    def get_model_name(self) -> str:
        """Return the model identifier."""
        pass
```

---

## S2-2: Claude Client

```python
# app/clients/claude_client.py
import anthropic
from app.clients.llm_client import LLMClient, ScoringResult
from app.config import settings

class ClaudeClient(LLMClient):
    def __init__(self):
        self.client = anthropic.AsyncAnthropic(api_key=settings.anthropic_api_key)
        self.model = settings.claude_model  # "claude-3-5-haiku-20241022"

    async def score_essay(
        self,
        content: str,
        prompt: str,
        few_shot_examples: list[dict]
    ) -> ScoringResult:
        messages = self._build_messages(content, few_shot_examples)

        response = await self.client.messages.create(
            model=self.model,
            max_tokens=1024,
            system=prompt,
            messages=messages
        )

        return self._parse_response(response.content[0].text)

    def _build_messages(self, content: str, examples: list[dict]) -> list[dict]:
        messages = []

        # Few-shot examples
        for ex in examples:
            messages.append({"role": "user", "content": f"다음 답안을 채점해주세요:\n\n{ex['content']}"})
            messages.append({"role": "assistant", "content": json.dumps(ex['result'], ensure_ascii=False)})

        # Actual request
        messages.append({"role": "user", "content": f"다음 답안을 채점해주세요:\n\n{content}"})

        return messages

    def _parse_response(self, response_text: str) -> ScoringResult:
        # Extract JSON from response
        try:
            result = json.loads(response_text)
            return ScoringResult(
                scores=result['scores'],
                feedback=result['feedback'],
                level_estimate=result.get('level_estimate', '')
            )
        except json.JSONDecodeError:
            # Try to extract JSON from markdown code block
            match = re.search(r'```json\s*(.*?)\s*```', response_text, re.DOTALL)
            if match:
                return json.loads(match.group(1))
            raise ValueError("Failed to parse LLM response")

    def get_model_name(self) -> str:
        return self.model
```

---

## S2-3: OpenAI Client

```python
# app/clients/openai_client.py
from openai import AsyncOpenAI
from app.clients.llm_client import LLMClient, ScoringResult
from app.config import settings

class OpenAIClient(LLMClient):
    def __init__(self):
        self.client = AsyncOpenAI(api_key=settings.openai_api_key)
        self.model = settings.openai_model  # "gpt-4o-mini"

    async def score_essay(
        self,
        content: str,
        prompt: str,
        few_shot_examples: list[dict]
    ) -> ScoringResult:
        messages = [{"role": "system", "content": prompt}]
        messages.extend(self._build_messages(content, few_shot_examples))

        response = await self.client.chat.completions.create(
            model=self.model,
            messages=messages,
            response_format={"type": "json_object"},
            max_tokens=1024
        )

        return self._parse_response(response.choices[0].message.content)

    def _build_messages(self, content: str, examples: list[dict]) -> list[dict]:
        messages = []
        for ex in examples:
            messages.append({"role": "user", "content": f"다음 답안을 채점해주세요:\n\n{ex['content']}"})
            messages.append({"role": "assistant", "content": json.dumps(ex['result'], ensure_ascii=False)})
        messages.append({"role": "user", "content": f"다음 답안을 채점해주세요:\n\n{content}"})
        return messages

    def _parse_response(self, response_text: str) -> ScoringResult:
        result = json.loads(response_text)
        return ScoringResult(
            scores=result['scores'],
            feedback=result['feedback'],
            level_estimate=result.get('level_estimate', '')
        )

    def get_model_name(self) -> str:
        return self.model
```

---

## S2-4: Prompt Builder

```python
# app/prompts/topik_scoring.py

SYSTEM_PROMPT = """당신은 TOPIK II 쓰기 54번 채점 전문가입니다.
공식 채점 기준에 따라 답안을 평가하고 구조화된 피드백을 제공합니다.

## 채점 기준

### 내용 및 과제 수행 (20점)
- 상(16-20): 주제에 맞게 자신의 생각을 풍부하게 전개함. 적절한 근거와 예시를 들어 주장을 뒷받침함.
- 중(10-15): 주제에 맞게 의견을 제시했으나 내용이 다소 부족하거나 근거가 충분하지 않음.
- 하(0-9): 주제에서 벗어나거나 내용이 매우 빈약함.

### 글의 전개 구조 (15점)
- 상(12-15): 서론-본론-결론 구조가 명확함. 문단 간 연결이 자연스럽고 논리적임.
- 중(7-11): 구조는 있으나 문단 연결이 다소 부자연스럽거나 논리 전개가 약함.
- 하(0-6): 구조가 불명확하거나 논리적 비약이 있음.

### 언어 사용 (15점)
- 상(12-15): 다양하고 적절한 어휘와 문법 사용. 고급 표현과 연결어미 활용.
- 중(7-11): 중급 수준의 어휘와 문법 사용. 일부 오류가 있으나 의미 전달에 지장 없음.
- 하(0-6): 어휘가 제한적이고 문법 오류가 빈번함.

## 출력 형식

반드시 다음 JSON 형식으로만 응답하세요:

```json
{
  "scores": {
    "content": <0-20>,
    "structure": <0-15>,
    "language": <0-15>,
    "total": <0-50>
  },
  "feedback": {
    "content": "<내용 평가 피드백>",
    "structure": "<구조 평가 피드백>",
    "language": "<언어 사용 피드백>",
    "overall": "<종합 피드백 및 개선 제안>"
  },
  "level_estimate": "<추정 TOPIK 등급: 3급/4급/5급/6급>"
}
```

## 채점 지침

1. 각 항목을 독립적으로 평가하세요.
2. 점수는 기준에 따라 객관적으로 부여하세요.
3. 피드백은 구체적이고 건설적으로 작성하세요.
4. 학습자가 개선할 수 있는 방향을 제시하세요.
5. 한국어로 피드백을 작성하세요."""


def build_prompt() -> str:
    return SYSTEM_PROMPT
```

---

## S2-5: Few-shot Data Loader

```python
# app/data/few_shot_loader.py
import json
from pathlib import Path
from functools import lru_cache

FEW_SHOT_PATH = Path(__file__).parent / "few_shot_examples.json"

@lru_cache(maxsize=1)
def load_few_shot_examples() -> list[dict]:
    """Load and cache few-shot examples."""
    with open(FEW_SHOT_PATH, "r", encoding="utf-8") as f:
        data = json.load(f)
    return data["examples"]


def get_examples_for_scoring(n: int = 3) -> list[dict]:
    """
    Get balanced examples for few-shot prompting.
    Returns 1 high, 1 medium, 1 low score example by default.
    """
    examples = load_few_shot_examples()

    # Group by level
    high = [e for e in examples if e["level"] == "고득점"]
    medium = [e for e in examples if e["level"] == "중간"]
    low = [e for e in examples if e["level"] == "저점수"]

    selected = []
    if high:
        selected.append(high[0])
    if medium:
        selected.append(medium[0])
    if low:
        selected.append(low[0])

    return selected[:n]
```

---

## S2-6: TOPIK Scorer Service

```python
# app/services/topik_scorer.py
from app.clients.llm_client import LLMClient, ScoringResult
from app.clients.claude_client import ClaudeClient
from app.clients.openai_client import OpenAIClient
from app.prompts.topik_scoring import build_prompt
from app.data.few_shot_loader import get_examples_for_scoring
from app.config import settings
import time

class TopikScorer:
    def __init__(self):
        self.client = self._get_client()
        self.prompt = build_prompt()

    def _get_client(self) -> LLMClient:
        if settings.llm_provider == "claude":
            return ClaudeClient()
        elif settings.llm_provider == "openai":
            return OpenAIClient()
        else:
            raise ValueError(f"Unknown LLM provider: {settings.llm_provider}")

    async def score(self, content: str) -> dict:
        """Score an essay and return full result with metadata."""
        examples = get_examples_for_scoring()

        start_time = time.time()
        result = await self.client.score_essay(content, self.prompt, examples)
        processing_time = int((time.time() - start_time) * 1000)

        return {
            "scores": result["scores"],
            "feedback": result["feedback"],
            "level_estimate": result["level_estimate"],
            "model": self.client.get_model_name(),
            "processing_time_ms": processing_time
        }
```

---

## S2-7: Worker Update

```python
# app/worker.py
import asyncio
import json
from app.services.topik_scorer import TopikScorer
from app.services.detector import DetectorService
from app.services.callback import CallbackService
from app.config import settings
import redis.asyncio as redis

class ScoringWorker:
    def __init__(self):
        self.redis = redis.from_url(settings.redis_url)
        self.scorer = TopikScorer()
        self.detector = DetectorService()
        self.callback = CallbackService()

    async def run(self):
        print("Scoring worker started...")
        while True:
            try:
                # Blocking pop from queue
                _, task_json = await self.redis.brpop("topik_scoring_queue")
                task = json.loads(task_json)

                await self.process_task(task)
            except Exception as e:
                print(f"Worker error: {e}")
                await asyncio.sleep(1)

    async def process_task(self, task: dict):
        submission_id = task["submission_id"]
        content = task["content"]

        try:
            # 1. AI Detection (local, fast)
            ai_detection = await self.detector.detect(content)

            # 2. LLM Scoring (API, slower)
            scoring_result = await self.scorer.score(content)

            # 3. Combine results
            result = {
                "submission_id": submission_id,
                "scores": scoring_result["scores"],
                "feedback": scoring_result["feedback"],
                "level_estimate": scoring_result["level_estimate"],
                "ai_detection": {
                    "score": ai_detection["score"],
                    "flagged": ai_detection["flagged"],
                    "model": ai_detection["model"]
                },
                "llm_model": scoring_result["model"],
                "processing_time_ms": scoring_result["processing_time_ms"],
                "status": "completed"
            }

            # 4. Send callback
            await self.callback.send(result)

        except Exception as e:
            # Send error callback
            await self.callback.send({
                "submission_id": submission_id,
                "status": "failed",
                "error": str(e)
            })


if __name__ == "__main__":
    worker = ScoringWorker()
    asyncio.run(worker.run())
```

---

## S2-10: Config Update

```python
# app/config.py
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    # Existing
    redis_url: str = "redis://localhost:6379"
    callback_base_url: str = "http://localhost:8080"
    callback_secret: str = ""

    # LLM Provider
    llm_provider: str = "claude"  # "claude" | "openai"

    # Claude
    anthropic_api_key: str = ""
    claude_model: str = "claude-3-5-haiku-20241022"

    # OpenAI
    openai_api_key: str = ""
    openai_model: str = "gpt-4o-mini"

    # AI Detection
    detector_model_path: str = "detector/chatgpt-detector-roberta-ko"

    class Config:
        env_file = ".env"

settings = Settings()
```

---

## LLM Response Validation

```python
# app/validators/scoring_validator.py
from pydantic import BaseModel, field_validator

class ScoresSchema(BaseModel):
    content: int
    structure: int
    language: int
    total: int

    @field_validator('content')
    def validate_content(cls, v):
        if not 0 <= v <= 20:
            raise ValueError('content score must be 0-20')
        return v

    @field_validator('structure', 'language')
    def validate_other_scores(cls, v):
        if not 0 <= v <= 15:
            raise ValueError('score must be 0-15')
        return v

    @field_validator('total')
    def validate_total(cls, v, values):
        expected = values.data.get('content', 0) + values.data.get('structure', 0) + values.data.get('language', 0)
        if v != expected:
            raise ValueError(f'total must equal sum of scores ({expected})')
        return v

class ScoringResponseSchema(BaseModel):
    scores: ScoresSchema
    feedback: dict[str, str]
    level_estimate: str
```

---

## Completion Criteria

- [ ] LLM Client 추상화 및 구현 (Claude, OpenAI)
- [ ] 프롬프트 빌더 구현
- [ ] Few-shot 데이터 로더 구현
- [ ] TOPIK 채점 서비스 구현
- [ ] Worker 업데이트
- [ ] AI 감지 통합
- [ ] Callback 페이로드 수정
- [ ] Config 업데이트
- [ ] 응답 검증 로직 구현
- [ ] 단위 테스트 작성

---

*Sprint 1 (API Server)와 병렬 진행 가능*
