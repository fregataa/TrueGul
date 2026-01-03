# Sprint 1: 인증 시스템

## 목표

회원가입, 로그인, 로그아웃 기능 구현

---

## 작업 목록

### Backend (API Server)

- [ ] 사용자 모델 및 레포지토리 구현
- [ ] 비밀번호 해싱 (bcrypt)
- [ ] JWT 토큰 생성/검증 유틸리티
- [ ] CSRF 토큰 생성/검증
- [ ] POST /api/auth/signup 구현
- [ ] POST /api/auth/login 구현 (HttpOnly Cookie 설정)
- [ ] POST /api/auth/logout 구현 (쿠키 삭제)
- [ ] 인증 미들웨어 구현
- [ ] 입력값 검증 (이메일 형식, 비밀번호 규칙)
- [ ] 에러 응답 형식 통일 (error_code, message)

### Frontend

- [ ] 인증 상태 관리 (Context 또는 Zustand)
- [ ] 회원가입 페이지 및 폼
- [ ] 로그인 페이지 및 폼
- [ ] 로그아웃 기능
- [ ] 인증 필요 페이지 보호 (미들웨어)
- [ ] CSRF 토큰 처리

### 테스트

- [ ] 인증 API 단위 테스트
- [ ] 인증 플로우 통합 테스트

---

## 완료 조건

- [ ] 회원가입 후 로그인 가능
- [ ] 로그인 시 HttpOnly Cookie로 JWT 저장
- [ ] 로그아웃 시 쿠키 삭제
- [ ] 미인증 사용자는 보호된 페이지 접근 불가
- [ ] 인증 API 테스트 통과
