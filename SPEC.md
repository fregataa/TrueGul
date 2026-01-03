# TrueGul 기획서

AI 없이 글쓰기 훈련을 돕는 서비스

---

## 개요

### 서비스 소개

TrueGul은 사용자가 AI 도움 없이 직접 글을 작성하고, 그 글이 AI로 작성되었는지 판별한 뒤 품질 피드백을 제공하는 서비스입니다.

### 목표

- 사용자의 순수 글쓰기 능력 향상
- AI 의존 없이 글 쓰는 습관 형성
- 글 종류별 맞춤 피드백 제공

---

## 로드맵

| 버전 | 목표 | 주요 내용 |
|------|------|-----------|
| v0 | MVP | 단일 모델, 기본 기능, 최소 로깅 |
| v1 | 모델 분리 | AI 탐지 / 피드백 모델 분리 |
| v2 | 모니터링 | ML 모니터링 + 서버 모니터링 |
| v3 | 타이핑 분석 | 타이핑 패턴 수집 및 분석 |

---

## v0 상세 명세

### 기능 요구사항

#### 1. 회원 관리

| 기능 | 설명 |
|------|------|
| 회원가입 | 이메일, 비밀번호 |
| 로그인 | JWT 기반 인증 |
| 로그아웃 | 토큰 무효화 |

#### 2. 글 작성

| 기능 | 설명 |
|------|------|
| 글 종류 선택 | 에세이, 자소서 |
| 글 작성 | 텍스트 에디터 |
| 임시 저장 | 작성 중 저장 |
| 글 제출 | 분석 요청 |
| 제출 이력 | 과거 제출 목록 조회 |

#### 3. AI 분석

| 기능 | 설명 |
|------|------|
| AI 탐지 | AI 작성 확률 점수 (0~100%) |
| 품질 피드백 | 글 종류별 맞춤 피드백 |
| Rate Limit | 사용자당 일일 제출 횟수 제한 |

### 비기능 요구사항

| 항목 | 요구사항 |
|------|----------|
| 응답시간 | AI 분석 결과 10초 이내 |
| 가용성 | 99% 이상 |
| 로깅 | 입력 텍스트, 모델 출력, 지연시간 기록 |
| 글 길이 제한 | 최대 2,000자 |

### 인증 방식

| 항목 | 설명 |
|------|------|
| 토큰 저장 | HttpOnly Cookie |
| 보안 | CSRF 토큰 사용 |
| 토큰 만료 | 1시간 |
| 로그아웃 | 클라이언트에서 쿠키 삭제 |

*Note: v0에서는 Refresh Token 및 서버사이드 토큰 무효화(블랙리스트) 미구현*

### API ↔ ML Server 통신

| 항목 | 설명 |
|------|------|
| 통신 방식 | 비동기 (Polling) |
| 제출 응답 | 202 Accepted + job_id 반환 |
| 결과 조회 | GET /api/writings/:id/analysis |
| 상태 값 | pending, processing, completed, failed |

### 에러 응답 형식

```json
{
  "error_code": "RATE_LIMIT_EXCEEDED",
  "message": "너무 많이 시도했습니다."
}
```

| 에러 코드 | 설명 |
|-----------|------|
| VALIDATION_ERROR | 입력값 검증 실패 |
| UNAUTHORIZED | 인증 필요 |
| RATE_LIMIT_EXCEEDED | 일일 제출 횟수 초과 |
| CONTENT_TOO_LONG | 글 길이 초과 (2,000자) |
| ANALYSIS_FAILED | AI 분석 실패 |
| NOT_FOUND | 리소스 없음 |

---

## 아키텍처

### 시스템 구성도

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Web UI    │────▶│  API Server │────▶│  Database   │
│  (Next.js)  │     │ (Go/Fargate)│     │ (RDS)       │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                           ▼ (비동기/Polling)
                    ┌─────────────┐     ┌─────────────┐
                    │  ML Server  │────▶│ OpenAI API  │
                    │  (Fargate)  │     │ (피드백)    │
                    └─────────────┘     └─────────────┘
```

### 기술 스택

| 구성요소 | 기술 |
|----------|------|
| Frontend | Next.js 15, React 19, TypeScript |
| API Server | Go |
| Database | PostgreSQL (AWS RDS) |
| ML Server | Python, FastAPI |
| 배포 | AWS (ECS/Fargate, RDS, S3 등) |

---

## API 설계

### 인증

| Method | Endpoint | 설명 |
|--------|----------|------|
| POST | /api/auth/signup | 회원가입 |
| POST | /api/auth/login | 로그인 |
| POST | /api/auth/logout | 로그아웃 |

### 글 관리

| Method | Endpoint | 설명 |
|--------|----------|------|
| GET | /api/writings | 글 목록 조회 |
| GET | /api/writings/:id | 글 상세 조회 |
| POST | /api/writings | 글 생성 (임시 저장) |
| PUT | /api/writings/:id | 글 수정 |
| DELETE | /api/writings/:id | 글 삭제 |
| POST | /api/writings/:id/submit | 글 제출 (분석 요청) |

### 분석 결과

| Method | Endpoint | 설명 |
|--------|----------|------|
| GET | /api/writings/:id/analysis | 분석 결과 조회 |

---

## 데이터 모델

### users

| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | UUID | PK |
| email | VARCHAR | 이메일 |
| password_hash | VARCHAR | 비밀번호 해시 |
| daily_submit_count | INTEGER | 당일 제출 횟수 (기본값: 0) |
| last_submit_date | DATE | 마지막 제출 날짜 |
| created_at | TIMESTAMP | 생성일시 |
| updated_at | TIMESTAMP | 수정일시 |

*Rate Limit 처리: 제출 시 `last_submit_date`가 오늘이 아니면 `daily_submit_count`를 0으로 리셋 후 증가. 별도 배치잡 없이 API 조회 시 계산.*

### writings

| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | UUID | PK |
| user_id | UUID | FK (users) |
| type | VARCHAR | 글 종류 (essay, cover_letter) |
| title | VARCHAR | 제목 |
| content | TEXT | 본문 |
| status | VARCHAR | 상태 (draft, submitted) |
| created_at | TIMESTAMP | 생성일시 |
| updated_at | TIMESTAMP | 수정일시 |
| submitted_at | TIMESTAMP | 제출일시 |

### analyses

| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | UUID | PK |
| writing_id | UUID | FK (writings) |
| status | VARCHAR | 상태 (pending, processing, completed, failed) |
| ai_score | FLOAT | AI 작성 확률 (0~1), nullable |
| feedback | TEXT | 품질 피드백, nullable |
| error_message | TEXT | 실패 시 에러 메시지, nullable |
| latency_ms | INTEGER | 분석 소요시간 |
| created_at | TIMESTAMP | 생성일시 |
| updated_at | TIMESTAMP | 수정일시 |

### analysis_logs (로깅용)

| 컬럼 | 타입 | 설명 |
|------|------|------|
| id | UUID | PK |
| analysis_id | UUID | FK (analyses) |
| input_text | TEXT | 입력 텍스트 |
| model_version | VARCHAR | 모델 버전 |
| raw_output | JSONB | 모델 원본 출력 |
| created_at | TIMESTAMP | 생성일시 |

---

## 글 종류별 피드백 기준

### 에세이

- 논리적 구조 (서론-본론-결론)
- 주장과 근거의 연결
- 문장 가독성
- 어휘 다양성

### 자소서

- 질문에 대한 답변 적합성
- 구체적 경험 서술
- 자기 강점 어필
- 진정성

---

## v1 ~ v3 개요

### v1: 모델 분리

- AI 탐지 모델: HuggingFace AI text detector (직접 배포 유지)
- 피드백 모델: LLM API (OpenAI) 또는 오픈소스 LLM 직접 배포 (선택)
- 모델별 독립 엔드포인트
- 모델 버전 관리

### v2: 모니터링

**ML 모니터링**
- 입력 데이터 드리프트
- 예측 분포 변화
- 모델별 추론 지연시간

**서버 모니터링**
- CPU, 메모리, 디스크
- 컨테이너 상태
- API 응답시간, 에러율, RPS

**도구**
- Prometheus + Grafana
- Evidently 또는 직접 구현

### v3: 타이핑 패턴 분석

- 키 입력 스트리밍 수집
- 분석 피처: 타이핑 속도, 멈춤 구간, 수정 패턴, 복붙 탐지
- 글쓰기 습관 리포트 제공

---

## 배포 계획 (v0)

### 인프라 구성 (AWS)

| 구성요소 | AWS 서비스 |
|----------|------------|
| Frontend | Vercel (또는 S3 + CloudFront) |
| API Server | ECS/Fargate |
| Database | RDS (PostgreSQL) |
| AI 탐지 서버 | ECS/Fargate (FastAPI + Docker) |
| 피드백 서버 | OpenAI API 호출 |
| Container Registry | ECR |

### 모델 배포 방식

MLOps 학습 목적으로 서버리스 GPU 서비스(Modal, Replicate) 대신 직접 배포 방식을 채택합니다.

**AI 탐지 모델 (직접 배포)**

```
HuggingFace 모델 → FastAPI 래핑 → Docker 이미지 → ECR → ECS/Fargate
```

- HuggingFace에서 AI text detector 모델 선정
- FastAPI로 추론 API 구현
- Docker 이미지로 패키징
- ECR에 이미지 푸시
- ECS/Fargate로 컨테이너 배포
- CPU 기반 모델 우선 선정 (비용 절감)

**피드백 모델 (API 호출)**

- v0: OpenAI API 호출 방식
- v1 이후: 오픈소스 LLM 직접 배포 검토 (선택사항)

### 직접 배포의 장점

- 컨테이너화 경험 (Docker, ECR, ECS)
- 모델 서빙 아키텍처 이해
- 로깅/모니터링 직접 구현 (CloudWatch)
- AWS 기반 MLOps 포트폴리오로 어필 가능

---

## 향후 고려사항

- 소셜 로그인 (Google, GitHub)
- 글 종류 확장 (논술, 블로그 등)
- 사용자별 글쓰기 성장 리포트
- 글 공유 기능
