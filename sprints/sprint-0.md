# Sprint 0: Planning & Setup

## Overview

| 항목 | 내용 |
|------|------|
| 목표 | 개발 환경 구성 및 기술 검증 |
| 선행 조건 | 없음 |
| 후속 Sprint | S1, S2, S3 (병렬 시작 가능) |

---

## Tasks

| ID | Task | 담당 | 상태 | 산출물 |
|----|------|------|------|--------|
| S0-1 | 모바일 프레임워크 선정 | - | DONE | Flutter 선정 |
| S0-2 | 모바일 프로젝트 초기화 | - | DONE | `/mobile` 디렉토리 |
| S0-3 | LLM API 선정 및 API Key 발급 | - | DONE | `.env.example`, `ml-server/app/config.py` |
| S0-4 | TOPIK 채점 프롬프트 초안 작성 | - | DONE | `ml-server/prompts/topik_scoring.md` |
| S0-5 | Few-shot 예시 데이터 준비 | - | DONE | `ml-server/data/few_shot_examples.json` |
| S0-6 | FCM/APNs 프로젝트 설정 | - | DEFERRED | → Sprint 3로 이동 |

---

## S0-1: 모바일 프레임워크 선정

### 비교

| 항목 | Flutter | React Native |
|------|---------|--------------|
| 언어 | Dart | JavaScript/TypeScript |
| 성능 | 높음 (네이티브 컴파일) | 중간 (JS Bridge) |
| UI 일관성 | 높음 (자체 렌더링) | 플랫폼별 차이 |
| 생태계 | 성장 중 | 성숙 |
| 학습 곡선 | Dart 학습 필요 | JS 경험 시 낮음 |

### 결정: **Flutter**
- [x] 높은 성능과 UI 일관성
- [x] 모바일 최적화
- [x] 단일 코드베이스로 iOS/Android 지원

---

## S0-3: LLM API 선정

### 후보

| 모델 | 예상 비용/회 | 장점 | 단점 |
|------|-------------|------|------|
| GPT-4o-mini | $0.002 | 저렴, 빠름 | 품질 검증 필요 |
| Claude 3.5 Haiku | $0.01 | 한국어 우수 | 비용 약간 높음 |
| GPT-4o | $0.03 | 품질 안정 | 비용 높음 |
| Claude 3.5 Sonnet | $0.04 | 최고 품질 | 비용 가장 높음 |

### 결정: **Claude + GPT 둘 다 테스트**
- [x] Claude 3.5 Haiku와 GPT-4o-mini 모두 설정
- [x] `scripts/test_llm_scoring.py`로 비교 테스트
- [ ] 테스트 결과 기반 최종 선정 (S2에서 진행)

---

## S0-4: TOPIK 채점 프롬프트 초안

### System Prompt 구조

```markdown
# 역할 정의
당신은 TOPIK II 쓰기 54번 채점 전문가입니다.

# 채점 기준
## 내용 및 과제 수행 (20점)
- 상(16-20): ...
- 중(10-15): ...
- 하(0-9): ...

## 글의 전개 구조 (15점)
...

## 언어 사용 (15점)
...

# Few-shot 예시
[예시 답안과 채점 결과]

# 출력 형식
JSON 스키마 정의
```

### 검증 방법
- 동일 답안 3회 채점 → 일관성 확인
- 공식 모범답안 채점 → 만점 근접 확인
- 의도적 오류 답안 → 적절한 감점 확인

---

## S0-5: Few-shot 예시 데이터

### 데이터 구성

| 점수대 | 예시 수 | 생성 방법 |
|--------|--------|----------|
| 고득점 (40-50) | 1개 | 공식 모범답안 |
| 중간 (25-35) | 2개 | 모범답안 의도적 열화 |
| 저점수 (10-20) | 2개 | 복합 감점 요인 적용 |

### 열화 기법

| 감점 영역 | 적용 기법 |
|----------|----------|
| 내용 | 근거 제거, 예시 부적절화 |
| 구성 | 결론 제거, 문단 미구분 |
| 표현 | 구어체 삽입, 문법 오류 추가, 어휘 반복 |

### 파일 형식

```json
{
  "examples": [
    {
      "id": "high-1",
      "level": "고득점",
      "content": "...",
      "scores": {
        "content": 18,
        "structure": 14,
        "language": 13,
        "total": 45
      },
      "feedback": {
        "content": "...",
        "structure": "...",
        "language": "...",
        "overall": "..."
      }
    }
  ]
}
```

---

## S0-6: Push Notification 설정

### Firebase (Android)
- [ ] Firebase 프로젝트 생성
- [ ] `google-services.json` 다운로드
- [ ] FCM API Key 발급
- [ ] 서버에 FCM 설정

### Apple (iOS)
- [ ] Apple Developer 계정
- [ ] Push Notification 인증서/키 생성
- [ ] APNs Key 다운로드
- [ ] 서버에 APNs 설정

---

## Completion Criteria

- [x] 모바일 프레임워크 결정 및 문서화 (Flutter)
- [x] 모바일 프로젝트 초기화 완료 (`/mobile`)
- [x] LLM API 선정 및 테스트 환경 구성 (Claude + GPT)
- [x] 프롬프트 초안 작성 (`ml-server/prompts/topik_scoring.md`)
- [x] Few-shot 데이터 5개 준비 (`ml-server/data/few_shot_examples.json`)
- [ ] ~~FCM/APNs 프로젝트 설정~~ → Sprint 3로 이동

---

**Sprint 0 Status: DONE** (2026-01-11)

*S1, S2, S3 병렬 진행 가능*
