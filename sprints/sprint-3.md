# Sprint 3: ML Server 및 AI 분석

## 목표

AI 탐지 모델 서빙 및 분석 파이프라인 구현

---

## 작업 목록

### ML Server

- [ ] HuggingFace AI text detector 모델 선정 및 테스트
- [ ] FastAPI 프로젝트 구조 설정
- [ ] POST /analyze 엔드포인트 구현 (AI 탐지)
- [ ] OpenAI API 연동 (품질 피드백 생성)
- [ ] 글 종류별 프롬프트 템플릿 작성 (에세이, 자소서)
- [ ] 응답 시간 측정 (latency_ms)
- [ ] 에러 처리 및 로깅
- [ ] Dockerfile 작성

### Backend (API Server)

- [ ] 분석 모델 및 레포지토리 구현
- [ ] POST /api/writings/:id/submit 구현
  - [ ] Rate Limit 체크 (daily_submit_count)
  - [ ] analyses 레코드 생성 (status: pending)
  - [ ] ML Server 호출 (비동기)
  - [ ] 202 Accepted 응답
- [ ] GET /api/writings/:id/analysis 구현 (Polling)
- [ ] ML Server 응답 처리 및 analyses 업데이트
- [ ] analysis_logs 저장 (입력 텍스트, 모델 출력, 버전)

### 테스트

- [ ] ML Server 단위 테스트 (모델 추론)
- [ ] 분석 API 통합 테스트
- [ ] Rate Limit 테스트

---

## 완료 조건

- [ ] ML Server가 텍스트를 받아 AI 탐지 점수 반환
- [ ] OpenAI API로 품질 피드백 생성
- [ ] 글 제출 시 202 Accepted 반환
- [ ] Polling으로 분석 결과 조회 가능
- [ ] Rate Limit 초과 시 에러 반환
- [ ] analysis_logs에 로그 저장
