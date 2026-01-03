# Sprint 0: 프로젝트 셋업

## 목표

개발 환경 및 인프라 기반 구축

---

## 작업 목록

### 프로젝트 초기화

- [ ] Frontend 프로젝트 생성 (Next.js 15, TypeScript, Tailwind, shadcn/ui)
- [ ] API Server 프로젝트 생성 (Go, 프로젝트 구조 설정)
- [ ] ML Server 프로젝트 생성 (Python, FastAPI)
- [ ] 모노레포 또는 멀티레포 구조 결정 및 설정

### 개발 환경

- [ ] Docker Compose 로컬 개발 환경 구성 (PostgreSQL, API Server, ML Server)
- [ ] 환경 변수 관리 방식 설정 (.env)
- [ ] Biome 설정 (Frontend)
- [ ] Go linter 설정 (API Server)
- [ ] Git hooks 설정 (Husky, pre-commit)

### 데이터베이스

- [ ] 데이터베이스 마이그레이션 도구 선정 및 설정 (golang-migrate 등)
- [ ] 초기 스키마 마이그레이션 작성 (users, writings, analyses, analysis_logs)

### CI/CD 기초

- [ ] GitHub Actions 워크플로우 기본 설정 (lint, test)

---

## 완료 조건

- [ ] 모든 프로젝트가 로컬에서 Docker Compose로 실행 가능
- [ ] 데이터베이스 마이그레이션이 정상 적용됨
- [ ] CI 파이프라인에서 lint가 통과함
