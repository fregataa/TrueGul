# Sprint 5: 배포 및 안정화

## 목표

AWS 배포 및 프로덕션 준비

---

## 작업 목록

### AWS 인프라

- [ ] VPC, 서브넷, 보안 그룹 설정
- [ ] RDS PostgreSQL 인스턴스 생성
- [ ] ECR 레포지토리 생성 (API Server, ML Server)
- [ ] ECS 클러스터 생성
- [ ] ECS 태스크 정의 (API Server, ML Server)
- [ ] ECS 서비스 생성 및 로드밸런서 연결
- [ ] 환경 변수 관리 (AWS Secrets Manager 또는 Parameter Store)

### CI/CD

- [ ] GitHub Actions 배포 워크플로우 (ECR 푸시, ECS 배포)
- [ ] 스테이징/프로덕션 환경 분리

### Frontend 배포

- [ ] Vercel 배포 설정
- [ ] 환경 변수 설정 (API URL)
- [ ] 도메인 연결 (선택)

### 모니터링 및 로깅

- [ ] CloudWatch 로그 그룹 설정
- [ ] 기본 알람 설정 (에러율, 응답시간)
- [ ] 헬스체크 엔드포인트

### 안정화

- [ ] 부하 테스트 (기본)
- [ ] 버그 수정 및 성능 개선
- [ ] 보안 점검 (HTTPS, CORS, 인증)

---

## 완료 조건

- [ ] 모든 서비스가 AWS에서 정상 동작
- [ ] CI/CD 파이프라인으로 자동 배포 가능
- [ ] CloudWatch에서 로그 확인 가능
- [ ] HTTPS로 서비스 접근 가능
- [ ] 기본 부하 테스트 통과
