# Version 3: 모델 최적화

## Overview

| 항목 | 내용 |
|------|------|
| 목표 | 채점 정확도 및 일관성 극대화 |
| 선행 조건 | v2 데이터 축적 완료 (1,000+ 채점 데이터) |
| 핵심 가치 | 실제 TOPIK 점수와 높은 상관관계 |

---

## Background

v2에서 축적된 데이터를 바탕으로:
1. AI 채점 모델을 최적화하고
2. 실제 TOPIK 점수와의 상관관계를 높입니다.

---

## Optimization Options

### Option A: Fine-tuning

**조건**: 채점 데이터 1,000+ 쌍 확보

| 항목 | 상세 |
|------|------|
| 학습 데이터 | 답안 + 채점 결과 (검증된 것만) |
| 기반 모델 | GPT-4o-mini 또는 Llama 계열 |
| 목표 | 도메인 특화 채점 모델 |

**장점**:
- 프롬프트 없이 직접 채점
- 응답 일관성 향상
- API 비용 절감 (자체 호스팅 시)

**단점**:
- 대량의 학습 데이터 필요
- 초기 투자 비용
- 모델 유지보수 필요

### Option B: Score Calibration

**조건**: 실제 TOPIK 점수 100+ 건 확보

| 항목 | 상세 |
|------|------|
| 데이터 | AI 점수 ↔ 실제 점수 쌍 |
| 방법 | 회귀 분석, 보정 함수 도출 |
| 적용 | AI 점수 후처리 |

**보정 함수 예시**:
```python
def calibrate_score(ai_score: int, question_type: str = "54") -> int:
    """AI 점수를 실제 TOPIK 점수에 맞게 보정"""
    # 선형 보정 예시
    # calibrated = a * ai_score + b
    # 계수는 회귀 분석으로 도출

    coefficients = {
        "54": {"a": 0.95, "b": 2.0},
    }

    coef = coefficients.get(question_type, {"a": 1.0, "b": 0.0})
    calibrated = coef["a"] * ai_score + coef["b"]

    return max(0, min(50, round(calibrated)))
```

**장점**:
- 상대적으로 적은 데이터로 가능
- 구현 간단
- 기존 시스템 변경 최소

**단점**:
- 근본적인 채점 품질 개선은 아님
- 지속적인 보정 필요

### Option C: RAG (Retrieval-Augmented Generation)

**조건**: 레벨별 표현 DB 구축

| 항목 | 상세 |
|------|------|
| 데이터 | TOPIK 레벨별 표현, 문법, 어휘 DB |
| 방법 | 답안 분석 시 관련 기준 검색하여 참조 |
| 적용 | 언어 사용 평가 정밀화 |

**구조**:
```
답안 입력
    ↓
[임베딩 생성]
    ↓
[표현 DB 검색] ←── 레벨별 어휘/문법 DB
    ↓
[관련 기준과 함께 LLM에 전달]
    ↓
정밀한 채점 결과
```

**표현 DB 예시**:
```json
{
  "level": "5급",
  "category": "연결어미",
  "expressions": [
    "-ㄴ/는 반면에",
    "-ㄴ/는 대신에",
    "-ㄴ/는 한편"
  ],
  "usage_examples": [...]
}
```

**장점**:
- 언어 평가 정밀도 향상
- 기준 업데이트 용이
- 피드백 구체성 향상

**단점**:
- DB 구축 비용
- 검색 정확도 의존
- 시스템 복잡도 증가

---

## Recommended Approach

### Phase 1: Score Calibration (우선)

v2 데이터로 즉시 적용 가능:

```python
# 1. 데이터 수집
# AI 점수 vs 실제 점수 100+ 쌍

# 2. 회귀 분석
from sklearn.linear_model import LinearRegression

model = LinearRegression()
model.fit(ai_scores, actual_scores)

# 3. 보정 함수 배포
calibration_a = model.coef_[0]
calibration_b = model.intercept_
```

### Phase 2: Prompt Enhancement

데이터 분석 기반 프롬프트 개선:

| 문제 패턴 | 개선 방안 |
|----------|----------|
| 고득점 과대평가 | 만점 기준 엄격화 |
| 문법 오류 과소평가 | 오류 유형별 감점 기준 명시 |
| 구어체 감지 부족 | 구어체 예시 추가 |

### Phase 3: Fine-tuning (장기)

데이터 충분 시 자체 모델 개발:

| 마일스톤 | 데이터 | 액션 |
|----------|--------|------|
| 1,000건 | 채점 데이터 | 파일럿 학습 테스트 |
| 5,000건 | 검증된 데이터 | 베타 모델 학습 |
| 10,000건 | 고품질 데이터 | 프로덕션 모델 배포 |

---

## Technical Implementation

### Calibration Service

```python
# app/services/calibration_service.py

class CalibrationService:
    def __init__(self):
        self.coefficients = self._load_coefficients()

    def _load_coefficients(self) -> dict:
        # DB 또는 설정에서 로드
        return {
            "54": {"a": 0.95, "b": 2.0, "updated_at": "2026-01-01"}
        }

    def calibrate(self, ai_score: int, question_type: str = "54") -> int:
        coef = self.coefficients.get(question_type, {"a": 1.0, "b": 0.0})
        calibrated = coef["a"] * ai_score + coef["b"]
        return max(0, min(50, round(calibrated)))

    def update_coefficients(self, question_type: str, new_coef: dict):
        # 새 계수 저장 (관리자 기능)
        pass
```

### Worker 수정

```python
# app/worker.py (수정)

async def process_task(self, task: dict):
    # ... 기존 코드 ...

    # 채점 결과에 보정 적용
    calibrated_total = self.calibration_service.calibrate(
        scoring_result["scores"]["total"]
    )

    result = {
        "submission_id": submission_id,
        "scores": {
            "content": scoring_result["scores"]["content"],
            "structure": scoring_result["scores"]["structure"],
            "language": scoring_result["scores"]["language"],
            "total": scoring_result["scores"]["total"],
            "calibrated_total": calibrated_total,  # 보정 점수 추가
        },
        # ...
    }
```

### Fine-tuning Pipeline (향후)

```python
# scripts/prepare_training_data.py

def prepare_training_data():
    """학습 데이터 준비"""
    # 1. 고품질 데이터 필터링
    # - 피드백 점수 4+
    # - 실제 점수와 AI 점수 오차 5점 이내

    # 2. 형식 변환 (OpenAI fine-tuning 형식)
    # {"messages": [
    #   {"role": "system", "content": "..."},
    #   {"role": "user", "content": "답안..."},
    #   {"role": "assistant", "content": "{scores, feedback}"}
    # ]}

    # 3. 학습/검증 분할 (80/20)
    pass
```

---

## Evaluation Metrics

### 보정 전후 비교

| 지표 | 보정 전 | 보정 후 목표 |
|------|---------|-------------|
| 상관계수 (r) | 0.65 | 0.80+ |
| MAE | 7점 | 4점 이하 |
| RMSE | 9점 | 5점 이하 |

### Fine-tuning 평가

| 지표 | 베이스라인 | 목표 |
|------|-----------|------|
| 점수 정확도 | 70% | 85%+ |
| 피드백 품질 | 3.8/5 | 4.5/5+ |
| 일관성 (재채점) | ±5점 | ±2점 |

---

## A/B Testing

새 모델 배포 전 A/B 테스트:

| 그룹 | 비율 | 모델 |
|------|------|------|
| A (Control) | 50% | 기존 모델 |
| B (Treatment) | 50% | 새 모델 |

**비교 지표**:
- 사용자 피드백 점수
- 재채점 요청률
- 이탈률

---

## Cost Considerations

### Fine-tuning 비용

| 항목 | 예상 비용 |
|------|----------|
| OpenAI Fine-tuning | $0.008/1K tokens (학습) |
| 10,000건 학습 | ~$50-100 |
| 월 추론 비용 | 기존 대비 20-30% 절감 가능 |

### 자체 호스팅 (Llama)

| 항목 | 예상 비용 |
|------|----------|
| GPU 인스턴스 (A10G) | $1-2/시간 |
| 월 운영 (24시간) | $720-1,440 |
| 대안: Spot Instance | $0.3-0.5/시간 |

**결론**: 초기에는 API Fine-tuning, 규모 증가 시 자체 호스팅 검토

---

## Success Criteria

| 지표 | 목표 |
|------|------|
| AI-실제 점수 상관계수 | 0.80+ |
| MAE | 4점 이하 |
| 사용자 피드백 점수 | 4.5/5.0+ |
| 재채점 일관성 | ±2점 이내 |

---

## Dependencies

- v2 데이터 축적 완료
  - 채점 데이터 1,000+ 건
  - 실제 점수 데이터 100+ 건
  - 사용자 피드백 데이터 충분

---

*v3 완료 후 v4 (확장) 진행*
