# Sprint 0: 프로젝트 셋업

## 목표

개발 환경 및 인프라 기반 구축

---

## 작업 목록

### 프로젝트 초기화

- [x] Frontend 프로젝트 생성 (Next.js 16, TypeScript, Tailwind v4, shadcn/ui)
- [x] API Server 프로젝트 생성 (Go, 프로젝트 구조 설정)
- [x] ML Server 프로젝트 생성 (Python, FastAPI)
- [x] 모노레포 또는 멀티레포 구조 결정 및 설정 → **모노레포 선택**

### 개발 환경

- [x] Docker Compose 로컬 개발 환경 구성 (PostgreSQL, API Server, ML Server)
- [x] 환경 변수 관리 방식 설정 (.env.example)
- [x] Biome 설정 (Frontend)
- [x] Go linter 설정 (API Server - golangci-lint)
- [x] Git hooks 설정 (Husky, pre-commit)

### 데이터베이스

- [x] 데이터베이스 마이그레이션 도구 선정 및 설정 (golang-migrate)
- [x] 초기 스키마 마이그레이션 작성 (users, writings, analyses, analysis_logs)

### CI/CD 기초

- [x] GitHub Actions 워크플로우 기본 설정 (lint, test)

---

## 완료 조건

- [x] 모든 프로젝트가 로컬에서 Docker Compose로 실행 가능
- [x] 데이터베이스 마이그레이션이 정상 적용됨
- [x] CI 파이프라인에서 lint가 통과함

---

## 구현 상세

### 프로젝트 구조

```
TrueGul/
├── .github/workflows/ci.yml
├── .husky/pre-commit
├── package.json (root - Husky)
├── docker-compose.yml
├── .env.example
│
├── frontend/
│   ├── biome.json
│   ├── components.json (shadcn/ui)
│   └── src/
│       ├── components/ui/button.tsx
│       └── lib/utils.ts
│
├── api-server/
│   ├── cmd/server/main.go
│   ├── internal/
│   ├── migrations/
│   │   ├── 000001_init_schema.up.sql
│   │   └── 000001_init_schema.down.sql
│   ├── .golangci.yml
│   ├── Dockerfile
│   └── Makefile
│
└── ml-server/
    ├── app/
    │   ├── main.py
    │   ├── config.py
    │   └── routes/analyze.py
    ├── tests/
    ├── pyproject.toml
    ├── ruff.toml
    └── Dockerfile
```

### 기술 스택 확정

| 컴포넌트 | 기술 |
|----------|------|
| Frontend | Next.js 16.1.1, React 19.2.3, TypeScript 5, Tailwind CSS v4, shadcn/ui |
| Frontend Linter | Biome |
| API Server | Go 1.21+, Gin, pgx |
| API Server Linter | golangci-lint |
| ML Server | Python 3.11+, FastAPI, Pydantic |
| ML Server Linter | Ruff |
| Database | PostgreSQL 15 |
| Migration | golang-migrate |
| Container | Docker, Docker Compose |
| CI/CD | GitHub Actions |
| Git Hooks | Husky |

### 로컬 실행 방법

```bash
# 환경 변수 설정
cp .env.example .env

# Docker Compose로 실행
docker-compose up -d

# pgAdmin 포함 실행 (선택사항)
docker-compose --profile tools up -d
```
