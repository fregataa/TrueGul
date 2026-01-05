# Sprint 5: 배포 및 안정화

## 목표

AWS 배포 및 프로덕션 준비

---

## 작업 목록

### AWS 인프라 (Terraform)

- [x] VPC, 서브넷, 보안 그룹 설정
- [x] RDS PostgreSQL 인스턴스 생성
- [x] ECR 레포지토리 생성 (API Server, ML Server)
- [x] ElastiCache Redis 생성
- [x] ECS 클러스터 생성
- [x] ECS 태스크 정의 (API Server, ML Server)
- [x] ECS 서비스 생성 및 로드밸런서 연결
- [x] 환경 변수 관리 (AWS Secrets Manager)
- [x] Auto Scaling 설정

### CI/CD

- [x] GitHub Actions 빌드/푸시 워크플로우 (ECR)
- [x] GitHub Actions 배포 워크플로우 (ECS)
- [x] 스테이징/프로덕션 환경 분리

### Frontend 배포

- [ ] Vercel 배포 설정 (수동 설정 필요)
- [ ] 환경 변수 설정 (API URL)
- [ ] 도메인 연결 (나중에)

### 모니터링 및 로깅

- [x] CloudWatch 로그 그룹 설정
- [ ] 기본 알람 설정 (에러율, 응답시간)
- [x] 헬스체크 엔드포인트 개선

### 안정화

- [ ] 부하 테스트 (기본)
- [ ] 버그 수정 및 성능 개선
- [ ] 보안 점검 (HTTPS, CORS, 인증)

---

## 구현 내역

### 신규 파일

| 디렉토리 | 설명 |
|----------|------|
| `infra/modules/vpc/` | VPC, 서브넷, 보안 그룹 모듈 |
| `infra/modules/ecr/` | ECR 레포지토리 모듈 |
| `infra/modules/rds/` | RDS PostgreSQL 모듈 |
| `infra/modules/elasticache/` | ElastiCache Redis 모듈 |
| `infra/modules/ecs/` | ECS 클러스터, 서비스, ALB 모듈 |
| `infra/environments/staging/` | Staging 환경 설정 |
| `infra/environments/production/` | Production 환경 설정 |
| `infra/README.md` | 인프라 배포 가이드 |
| `.github/workflows/build-push.yml` | ECR 이미지 빌드/푸시 |
| `.github/workflows/deploy.yml` | ECS 배포 워크플로우 |

### 수정된 파일

| 파일 | 변경 내용 |
|------|----------|
| `api-server/internal/handler/health.go` | Health check 핸들러 (DB, Redis 상태 포함) |
| `api-server/internal/mq/redis.go` | Client() getter 추가 |
| `api-server/cmd/server/main.go` | Health handler 통합 |
| `ml-server/app/main.py` | Health check 개선 (모델, OpenAI 상태) |
| `ml-server/app/services/detector.py` | is_loaded() 메서드 추가 |
| `ml-server/app/services/feedback.py` | is_initialized() 메서드 추가 |

---

## 배포 순서

1. **Terraform State Backend 생성** (S3 + DynamoDB)
2. **ECR 레포지토리 생성** (Staging Terraform)
3. **Docker 이미지 빌드/푸시** (GitHub Actions 또는 수동)
4. **AWS Secrets Manager 시크릿 등록**
5. **VPC, RDS, Redis 생성** (Terraform)
6. **ECS 클러스터/서비스 생성** (Terraform)
7. **DB 마이그레이션** (ECS Exec)
8. **Vercel 프로젝트 생성**
9. **E2E 테스트**

---

## 완료 조건

- [x] Terraform 인프라 코드 완성
- [x] GitHub Actions CI/CD 파이프라인 완성
- [x] Health check 엔드포인트 개선
- [ ] 모든 서비스가 AWS에서 정상 동작 (배포 후 확인)
- [ ] CloudWatch에서 로그 확인 가능 (배포 후 확인)
- [ ] HTTPS로 서비스 접근 가능 (도메인 연결 후)
- [ ] 기본 부하 테스트 통과

---

## AWS 예상 비용 (월간)

| 서비스 | Staging | Production | 합계 |
|--------|---------|------------|------|
| ECS Fargate | ~$15 | ~$30 | ~$45 |
| RDS | ~$15 | ~$30 | ~$45 |
| ElastiCache | ~$12 | ~$12 | ~$24 |
| NAT Gateway | ~$32 | ~$32 | ~$64 |
| ALB | ~$16 | ~$16 | ~$32 |
| ECR | ~$1 | ~$1 | ~$2 |
| **합계** | **~$91** | **~$121** | **~$212** |
