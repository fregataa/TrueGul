#!/usr/bin/env python3
"""
TOPIK 채점 프롬프트 테스트 스크립트

테스트 항목:
1. Claude vs GPT 비교
2. 일관성 테스트 (동일 답안 3회 채점)
3. 점수 범위 검증 (고득점 → 40+, 저득점 → 25-)

사용법:
    cd ml-server
    python scripts/test_llm_scoring.py --provider anthropic
    python scripts/test_llm_scoring.py --provider openai
    python scripts/test_llm_scoring.py --provider both
"""

import argparse
import json
import os
import sys
from pathlib import Path
from typing import Any

# Add project root to path
project_root = Path(__file__).parent.parent
sys.path.insert(0, str(project_root))

try:
    import anthropic
    import openai
except ImportError as e:
    print(f"Required packages not installed: {e}")
    print("Run: pip install anthropic openai")
    sys.exit(1)


def load_prompt() -> str:
    """Load TOPIK scoring prompt from markdown file."""
    prompt_path = project_root / "prompts" / "topik_scoring.md"
    with open(prompt_path, encoding="utf-8") as f:
        content = f.read()

    # Extract system prompt section
    # For simplicity, we'll construct a prompt from the markdown
    return content


def load_few_shot_examples() -> list[dict]:
    """Load few-shot examples from JSON file."""
    data_path = project_root / "data" / "few_shot_examples.json"
    with open(data_path, encoding="utf-8") as f:
        data = json.load(f)
    return data["examples"]


def build_system_prompt(prompt_content: str, examples: list[dict]) -> str:
    """Build system prompt with few-shot examples."""

    system_prompt = """당신은 TOPIK II 쓰기 54번 채점 전문가입니다.

TOPIK(한국어능력시험) 공식 채점 기준에 따라 600-700자 분량의 의견 제시형 논술 답안을 평가합니다.
객관적이고 일관된 채점을 수행하며, 각 항목별로 구체적인 피드백을 제공합니다.

## 채점 기준 (총 50점)

### 1. 내용 및 과제 수행 (20점)
- 상(16-20): 주제를 정확히 이해하고 풍부한 내용을 논리적으로 전개
- 중(10-15): 주제에 맞으나 근거가 다소 부족
- 하(0-9): 주제 이탈 또는 내용 빈약

### 2. 글의 전개 구조 (15점)
- 상(12-15): 서론-본론-결론 명확, 연결어 효과적 사용
- 중(7-11): 구조는 있으나 연결이 부자연스러움
- 하(0-6): 구조 불명확, 논리 비약

### 3. 언어 사용 (15점)
- 상(12-15): 다양한 어휘/문법, 고급 표현 사용
- 중(7-11): 중급 수준, 일부 오류
- 하(0-6): 제한적 어휘, 빈번한 문법 오류, 구어체 사용

## 출력 형식
반드시 아래 JSON 형식으로만 응답하세요:

```json
{
  "scores": {
    "content": <0-20>,
    "structure": <0-15>,
    "language": <0-15>,
    "total": <0-50>
  },
  "feedback": {
    "content": "<내용 피드백>",
    "structure": "<구성 피드백>",
    "language": "<언어 피드백>",
    "overall": "<종합 피드백>"
  },
  "estimated_level": "<1급~6급>"
}
```

## 채점 예시
"""

    # Add 2 examples (high and low score) for few-shot
    for ex in examples[:2]:  # Use first 2 examples
        system_prompt += f"""
### 예시 - {ex['level']} ({ex['scores']['total']}점)
문제: {ex['topic']}
답안: {ex['content'][:200]}...

채점 결과:
```json
{json.dumps({"scores": ex["scores"], "feedback": ex["feedback"], "estimated_level": ex["estimated_level"]}, ensure_ascii=False, indent=2)}
```
"""

    return system_prompt


def score_with_anthropic(
    system_prompt: str,
    topic: str,
    answer: str,
    api_key: str,
    model: str = "claude-3-5-haiku-20241022"
) -> dict[str, Any]:
    """Score an essay using Claude API."""
    client = anthropic.Anthropic(api_key=api_key)

    user_message = f"""다음 TOPIK II 54번 답안을 채점해주세요.

**문제:**
{topic}

**답안:**
{answer}

JSON 형식으로만 응답해주세요."""

    response = client.messages.create(
        model=model,
        max_tokens=1024,
        system=system_prompt,
        messages=[{"role": "user", "content": user_message}]
    )

    # Parse JSON from response
    response_text = response.content[0].text

    # Extract JSON from response (handle markdown code blocks)
    if "```json" in response_text:
        json_str = response_text.split("```json")[1].split("```")[0].strip()
    elif "```" in response_text:
        json_str = response_text.split("```")[1].split("```")[0].strip()
    else:
        json_str = response_text.strip()

    return json.loads(json_str)


def score_with_openai(
    system_prompt: str,
    topic: str,
    answer: str,
    api_key: str,
    model: str = "gpt-4o-mini"
) -> dict[str, Any]:
    """Score an essay using OpenAI API."""
    client = openai.OpenAI(api_key=api_key)

    user_message = f"""다음 TOPIK II 54번 답안을 채점해주세요.

**문제:**
{topic}

**답안:**
{answer}

JSON 형식으로만 응답해주세요."""

    response = client.chat.completions.create(
        model=model,
        max_tokens=1024,
        messages=[
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_message}
        ],
        response_format={"type": "json_object"}
    )

    response_text = response.choices[0].message.content
    return json.loads(response_text)


def test_consistency(
    score_func,
    system_prompt: str,
    topic: str,
    answer: str,
    api_key: str,
    num_trials: int = 3
) -> tuple[list[int], float]:
    """Test scoring consistency by running multiple times."""
    scores = []
    for i in range(num_trials):
        try:
            result = score_func(system_prompt, topic, answer, api_key)
            scores.append(result["scores"]["total"])
            print(f"  Trial {i+1}: {result['scores']['total']}점")
        except Exception as e:
            print(f"  Trial {i+1}: Error - {e}")

    if len(scores) >= 2:
        variance = max(scores) - min(scores)
        return scores, variance
    return scores, float('inf')


def run_tests(provider: str):
    """Run all tests for the specified provider."""
    print(f"\n{'='*60}")
    print(f"TOPIK 채점 프롬프트 테스트 - Provider: {provider}")
    print(f"{'='*60}\n")

    # Load data
    prompt_content = load_prompt()
    examples = load_few_shot_examples()
    system_prompt = build_system_prompt(prompt_content, examples)

    # Get API keys
    anthropic_key = os.getenv("ANTHROPIC_API_KEY", "")
    openai_key = os.getenv("OPENAI_API_KEY", "")

    providers_to_test = []
    if provider in ["anthropic", "both"]:
        if anthropic_key:
            providers_to_test.append(("anthropic", score_with_anthropic, anthropic_key))
        else:
            print("Warning: ANTHROPIC_API_KEY not set")

    if provider in ["openai", "both"]:
        if openai_key:
            providers_to_test.append(("openai", score_with_openai, openai_key))
        else:
            print("Warning: OPENAI_API_KEY not set")

    if not providers_to_test:
        print("Error: No API keys configured. Set ANTHROPIC_API_KEY or OPENAI_API_KEY")
        return

    # Test cases
    test_cases = [
        ("high-1", "고득점 검증", 40),  # Expected 40+
        ("low-1", "저득점 검증", 25),   # Expected 25-
    ]

    results = {}

    for name, score_func, api_key in providers_to_test:
        print(f"\n--- {name.upper()} 테스트 ---\n")
        results[name] = {"tests": [], "consistency": None}

        # Score validation tests
        for example_id, test_name, threshold in test_cases:
            example = next((ex for ex in examples if ex["id"] == example_id), None)
            if not example:
                continue

            print(f"[{test_name}] {example_id}")
            try:
                result = score_func(system_prompt, example["topic"], example["content"], api_key)
                score = result["scores"]["total"]
                expected = example["scores"]["total"]

                if example_id.startswith("high"):
                    passed = score >= threshold
                    condition = f">= {threshold}"
                else:
                    passed = score <= threshold
                    condition = f"<= {threshold}"

                status = "PASS" if passed else "FAIL"
                print(f"  결과: {score}점 (기대: {expected}점, 조건: {condition}) [{status}]")

                results[name]["tests"].append({
                    "test": test_name,
                    "score": score,
                    "expected": expected,
                    "passed": passed
                })

            except Exception as e:
                print(f"  Error: {e}")
                results[name]["tests"].append({
                    "test": test_name,
                    "error": str(e)
                })

        # Consistency test
        print(f"\n[일관성 테스트] 동일 답안 3회 채점")
        test_example = examples[1]  # mid-1
        scores, variance = test_consistency(
            score_func,
            system_prompt,
            test_example["topic"],
            test_example["content"],
            api_key
        )

        consistency_passed = variance <= 6  # ±3점 = 최대 6점 차이
        status = "PASS" if consistency_passed else "FAIL"
        print(f"  점수 범위: {min(scores)}-{max(scores)} (편차: {variance}점) [{status}]")

        results[name]["consistency"] = {
            "scores": scores,
            "variance": variance,
            "passed": consistency_passed
        }

    # Summary
    print(f"\n{'='*60}")
    print("테스트 요약")
    print(f"{'='*60}")

    for name, data in results.items():
        passed = sum(1 for t in data["tests"] if t.get("passed"))
        total = len(data["tests"])
        consistency = "PASS" if data["consistency"]["passed"] else "FAIL"
        print(f"\n{name.upper()}:")
        print(f"  점수 검증: {passed}/{total} 통과")
        print(f"  일관성: {consistency} (편차 {data['consistency']['variance']}점)")


def main():
    parser = argparse.ArgumentParser(description="Test TOPIK scoring prompts")
    parser.add_argument(
        "--provider",
        choices=["anthropic", "openai", "both"],
        default="both",
        help="LLM provider to test"
    )
    args = parser.parse_args()

    run_tests(args.provider)


if __name__ == "__main__":
    main()
