# Sprint 5: Deployment

## Overview

| 항목 | 내용 |
|------|------|
| 목표 | 프로덕션 배포 및 베타 테스트 시작 |
| 선행 조건 | Sprint 4 완료 (QA 통과) |
| 결과 | v0 클로즈드 베타 런칭 |

---

## Tasks

| ID | Task | 상태 | 우선순위 |
|----|------|------|----------|
| S5-1 | AWS 인프라 재시작 | TODO | P0 |
| S5-2 | 환경변수 설정 | TODO | P0 |
| S5-3 | DB 마이그레이션 실행 | TODO | P0 |
| S5-4 | API Server 배포 | TODO | P0 |
| S5-5 | ML Server 배포 | TODO | P0 |
| S5-6 | 앱 빌드 (iOS) | TODO | P0 |
| S5-7 | 앱 빌드 (Android) | TODO | P0 |
| S5-8 | 클로즈드 베타 테스터 모집 | TODO | P1 |
| S5-9 | 모니터링 설정 | TODO | P1 |
| S5-10 | 비용 모니터링 | TODO | P1 |

---

## S5-1: AWS Infrastructure Restart

### 현재 비활성화된 리소스

| 리소스 | 상태 | 재시작 방법 |
|--------|------|------------|
| ECS (API Server) | 중지됨 | desired_count 복원 |
| ECS (ML Server) | 중지됨 | desired_count 복원 |
| RDS | 중지됨 | terraform apply |
| ElastiCache | 중지됨 | terraform apply |

### 재시작 절차

```bash
cd infra/environments/production

# 1. deletion_protection 다시 활성화
# main.tf에서 deletion_protection = true 확인

# 2. terraform apply
terraform init
terraform plan
terraform apply

# 3. 상태 확인
aws ecs describe-services --cluster truegul-production --services api-server ml-server
```

### 예상 비용 (월)

| 리소스 | 예상 비용 |
|--------|----------|
| ECS Fargate (API) | $30-50 |
| ECS Fargate (ML) | $40-60 |
| RDS (db.t3.small) | $25-30 |
| ElastiCache (cache.t3.micro) | $12-15 |
| NAT Gateway | $32-45 |
| ALB | $16-25 |
| **Total** | **$155-225** |

---

## S5-2: Environment Variables

### API Server

```bash
# AWS Secrets Manager에 설정
aws secretsmanager create-secret \
  --name truegul/production/jwt \
  --secret-string "your-jwt-secret"

aws secretsmanager create-secret \
  --name truegul/production/callback \
  --secret-string "your-callback-secret"
```

### ML Server

| 변수 | 값 | 비고 |
|------|-----|------|
| `LLM_PROVIDER` | `claude` | 또는 `openai` |
| `ANTHROPIC_API_KEY` | `sk-ant-...` | Secrets Manager |
| `CLAUDE_MODEL` | `claude-3-5-haiku-20241022` | |
| `CALLBACK_BASE_URL` | `http://api-server:8080` | 내부 통신 |
| `CALLBACK_SECRET` | `...` | Secrets Manager |

### ECS Task Definition 업데이트

```hcl
# infra/modules/ecs/main.tf
# ML Server Task Definition에 추가

secrets = [
  {
    name      = "ANTHROPIC_API_KEY"
    valueFrom = "arn:aws:secretsmanager:${var.aws_region}:${var.aws_account_id}:secret:${var.project}/${var.environment}/anthropic"
  }
]

environment = [
  { name = "LLM_PROVIDER", value = "claude" },
  { name = "CLAUDE_MODEL", value = "claude-3-5-haiku-20241022" },
  # ...
]
```

---

## S5-3: Database Migration

### 새 마이그레이션 파일

```
migrations/
├── 001_create_users.sql (기존)
├── 002_create_writings.sql (기존)
├── 003_create_tasks.sql (기존)
├── 004_create_submissions.sql (신규)
├── 005_create_scoring_results.sql (신규)
└── 006_create_push_tokens.sql (신규)
```

### 마이그레이션 실행

```bash
# ECS Task로 마이그레이션 실행
aws ecs run-task \
  --cluster truegul-production \
  --task-definition truegul-production-migrate \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx],securityGroups=[sg-xxx],assignPublicIp=DISABLED}"

# 로그 확인
aws logs tail /ecs/truegul-production/migrate --follow
```

---

## S5-4: API Server Deployment

### Docker Build & Push

```bash
cd api-server

# Build
docker build -t truegul-api-server:v0.1.0 .

# Tag
docker tag truegul-api-server:v0.1.0 \
  ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-2.amazonaws.com/truegul-api-server:v0.1.0

# Push
aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-2.amazonaws.com

docker push ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-2.amazonaws.com/truegul-api-server:v0.1.0
```

### ECS Service Update

```bash
aws ecs update-service \
  --cluster truegul-production \
  --service api-server \
  --force-new-deployment
```

---

## S5-5: ML Server Deployment

### Docker Build & Push

```bash
cd ml-server

# Build (requirements.txt에 anthropic, openai 추가됨)
docker build -t truegul-ml-server:v0.1.0 .

# Tag & Push
docker tag truegul-ml-server:v0.1.0 \
  ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-2.amazonaws.com/truegul-ml-server:v0.1.0

docker push ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-2.amazonaws.com/truegul-ml-server:v0.1.0
```

### ECS Service Update

```bash
aws ecs update-service \
  --cluster truegul-production \
  --service ml-server \
  --force-new-deployment
```

---

## S5-6: iOS App Build

### Prerequisites

- [ ] Apple Developer 계정
- [ ] App Store Connect 앱 등록
- [ ] Provisioning Profile 생성
- [ ] Push Notification 인증서

### Build Steps

```bash
cd mobile

# Dependencies
flutter pub get

# Build for release
flutter build ipa --release

# 또는 Xcode에서
open ios/Runner.xcworkspace
# Product > Archive > Distribute App > App Store Connect
```

### TestFlight 배포

1. App Store Connect에서 앱 선택
2. TestFlight 탭
3. 빌드 업로드 확인
4. 테스터 그룹 추가
5. 테스터 초대 발송

---

## S5-7: Android App Build

### Prerequisites

- [ ] Google Play Console 앱 등록
- [ ] Keystore 생성
- [ ] Firebase 프로젝트 연결

### Build Steps

```bash
cd mobile

# Dependencies
flutter pub get

# Build AAB for Play Store
flutter build appbundle --release

# 생성 파일: build/app/outputs/bundle/release/app-release.aab
```

### Internal Testing 배포

1. Google Play Console에서 앱 선택
2. Release > Testing > Internal testing
3. Create new release
4. AAB 파일 업로드
5. 테스터 이메일 추가
6. Start rollout

---

## S5-8: Closed Beta Testers

### 모집 채널

| 채널 | 방법 |
|------|------|
| 지인 | 직접 연락 |
| TOPIK 커뮤니티 | 베타 테스터 모집 글 |
| 한국어 학습 커뮤니티 | Reddit, Discord 등 |

### 테스터 요건

- TOPIK 준비 중인 학습자
- iOS 또는 Android 기기 보유
- 피드백 제공 의향

### 모집 목표

| 단계 | 인원 | 목적 |
|------|------|------|
| 1차 (1주) | 5-10명 | 치명적 버그 발견 |
| 2차 (2주) | 20-30명 | 피드백 수집 |

### 피드백 수집

- 인앱 피드백 버튼
- 구글 폼 설문
- 1:1 인터뷰 (선택)

---

## S5-9: Monitoring Setup

### CloudWatch Dashboard

| 지표 | 알람 조건 |
|------|----------|
| ECS CPU 사용률 | > 80% |
| ECS 메모리 사용률 | > 80% |
| ALB 5xx 에러율 | > 1% |
| API 응답 시간 | > 1초 |
| RDS 연결 수 | > 80% |

### Log Groups

```
/ecs/truegul-production/api-server
/ecs/truegul-production/ml-server
/ecs/truegul-production/migrate
```

### Alerts

```bash
# SNS Topic 생성
aws sns create-topic --name truegul-production-alerts

# 이메일 구독
aws sns subscribe \
  --topic-arn arn:aws:sns:ap-northeast-2:xxx:truegul-production-alerts \
  --protocol email \
  --notification-endpoint your-email@example.com
```

---

## S5-10: Cost Monitoring

### LLM API 비용 추적

| 항목 | 추적 방법 |
|------|----------|
| 일별 요청 수 | CloudWatch 커스텀 메트릭 |
| 토큰 사용량 | ML Server 로그 |
| 비용 추정 | 요청 수 × 예상 단가 |

### 비용 알람 설정

```bash
# AWS Budgets 설정
aws budgets create-budget \
  --account-id ${AWS_ACCOUNT_ID} \
  --budget '{
    "BudgetName": "TrueGul-Monthly",
    "BudgetLimit": {"Amount": "300", "Unit": "USD"},
    "TimeUnit": "MONTHLY",
    "BudgetType": "COST"
  }' \
  --notifications-with-subscribers '[{
    "Notification": {
      "NotificationType": "ACTUAL",
      "ComparisonOperator": "GREATER_THAN",
      "Threshold": 80
    },
    "Subscribers": [{
      "SubscriptionType": "EMAIL",
      "Address": "your-email@example.com"
    }]
  }]'
```

### 월간 비용 리포트

| 항목 | 예산 | 실제 |
|------|------|------|
| AWS 인프라 | $200 | - |
| LLM API | $100 | - |
| **Total** | **$300** | - |

---

## Deployment Checklist

### Pre-deployment

- [ ] 모든 테스트 통과 (Sprint 4)
- [ ] 환경변수 준비
- [ ] Secrets Manager 설정
- [ ] Docker 이미지 빌드

### Deployment

- [ ] AWS 인프라 재시작
- [ ] DB 마이그레이션 실행
- [ ] API Server 배포
- [ ] ML Server 배포
- [ ] 헬스체크 확인
- [ ] E2E 테스트 (프로덕션)

### Post-deployment

- [ ] 모니터링 대시보드 확인
- [ ] 알람 설정 확인
- [ ] iOS TestFlight 배포
- [ ] Android Internal Testing 배포
- [ ] 테스터 초대 발송

---

## Rollback Plan

### API/ML Server

```bash
# 이전 버전으로 롤백
aws ecs update-service \
  --cluster truegul-production \
  --service api-server \
  --task-definition truegul-production-api-server:PREVIOUS_VERSION
```

### Database

```bash
# 마이그레이션 롤백 (down 스크립트 필요)
./migrate down
```

### Mobile App

- TestFlight/Internal Testing에서 이전 빌드 활성화
- 또는 새 수정 버전 빠른 배포

---

## Completion Criteria

- [ ] AWS 인프라 정상 가동
- [ ] API/ML Server 배포 완료
- [ ] 프로덕션 E2E 테스트 통과
- [ ] iOS TestFlight 배포 완료
- [ ] Android Internal Testing 배포 완료
- [ ] 모니터링 및 알람 설정 완료
- [ ] 최소 5명 베타 테스터 참여

---

## v0 Launch Summary

| 항목 | 상태 |
|------|------|
| 버전 | v0.1.0 |
| 플랫폼 | iOS (TestFlight), Android (Internal Testing) |
| 기능 | TOPIK 54번 채점, AI 감지 |
| 사용자 | 클로즈드 베타 (초대 전용) |

---

*v0 런칭 완료 후 v1 (OCR) 개발 시작*
