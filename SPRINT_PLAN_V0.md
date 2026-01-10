# TrueGul v0 Sprint Plan - Index

## Overview

| 항목 | 내용 |
|------|------|
| 목표 | TOPIK II 54번 AI 채점 MVP |
| 범위 | 텍스트 입력 → LLM 채점 → 결과 확인 |
| 플랫폼 | Mobile (Flutter 또는 React Native) |

---

## Sprint Documents

| Sprint | 목표 | 문서 |
|--------|------|------|
| **S0** | Planning & Setup | [sprint-0.md](./sprints/sprint-0.md) |
| **S1** | Backend - API Server | [sprint-1.md](./sprints/sprint-1.md) |
| **S2** | Backend - ML Server | [sprint-2.md](./sprints/sprint-2.md) |
| **S3** | Mobile App | [sprint-3.md](./sprints/sprint-3.md) |
| **S4** | Integration & QA | [sprint-4.md](./sprints/sprint-4.md) |
| **S5** | Deployment | [sprint-5.md](./sprints/sprint-5.md) |

---

## Dependencies

```
S0 (Planning)
 │
 ├──▶ S1 (API Server) ──┐
 │                      │
 ├──▶ S2 (ML Server) ───┼──▶ S4 (Integration) ──▶ S5 (Deployment)
 │                      │
 └──▶ S3 (Mobile App) ──┘
```

| 관계 | 설명 |
|------|------|
| S0 → S1, S2, S3 | S0 완료 후 병렬 시작 가능 |
| S1 ↔ S2 | API 스키마 공유, 병렬 진행 |
| S3 → S1 | API 완료 후 연동 테스트 |
| S4 → S1, S2, S3 | 모든 개발 완료 후 통합 |

---

## Version Documents

추후 버전 계획은 별도 문서로 관리:

| Version | 목표 | 문서 |
|---------|------|------|
| **v1** | OCR 도입 | [v1-ocr.md](./versions/v1-ocr.md) |
| **v2** | 데이터 축적 | [v2-data-collection.md](./versions/v2-data-collection.md) |
| **v3** | 모델 최적화 | [v3-model-optimization.md](./versions/v3-model-optimization.md) |
| **v4+** | 확장 | [v4-expansion.md](./versions/v4-expansion.md) |

---

## Quick Summary

### v0 핵심 기능
- 회원가입/로그인
- TOPIK 54번 답안 텍스트 입력
- AI 채점 (LLM API)
- AI 작성 감지 (RoBERTa)
- Push 알림으로 결과 전달
- 채점 기록 조회

### 기술 스택
- Mobile: Flutter 또는 React Native
- API Server: Go/Gin
- ML Server: FastAPI/Python
- LLM: Claude 3.5 Haiku 또는 GPT-4o-mini
- DB: PostgreSQL
- Queue: Redis
- Push: FCM/APNs
- Infra: AWS ECS

### 예상 비용 (월)
- AWS 인프라: $150-225
- LLM API: $100-200
- **Total**: ~$250-425

---

*작성일: 2026년 1월 9일*
