# Sprint 1: 인증 시스템

## 목표

회원가입, 로그인, 로그아웃 기능 구현

---

## 작업 목록

### Backend (API Server)

- [x] 사용자 모델 및 레포지토리 구현
- [x] 비밀번호 해싱 (bcrypt)
- [x] JWT 토큰 생성/검증 유틸리티
- [x] CSRF 토큰 생성/검증
- [x] POST /api/auth/signup 구현
- [x] POST /api/auth/login 구현 (HttpOnly Cookie 설정)
- [x] POST /api/auth/logout 구현 (쿠키 삭제)
- [x] 인증 미들웨어 구현
- [x] 입력값 검증 (이메일 형식, 비밀번호 규칙)
- [x] 에러 응답 형식 통일 (error_code, message)

### Frontend

- [x] 인증 상태 관리 (Context 또는 Zustand)
- [x] 회원가입 페이지 및 폼
- [x] 로그인 페이지 및 폼
- [x] 로그아웃 기능
- [x] 인증 필요 페이지 보호 (미들웨어)
- [x] CSRF 토큰 처리

### 테스트

- [ ] 인증 API 단위 테스트
- [ ] 인증 플로우 통합 테스트

---

## 완료 조건

- [x] 회원가입 후 로그인 가능
- [x] 로그인 시 HttpOnly Cookie로 JWT 저장
- [x] 로그아웃 시 쿠키 삭제
- [x] 미인증 사용자는 보호된 페이지 접근 불가
- [ ] 인증 API 테스트 통과

---

## 구현 상세

### Backend 구조

```
api-server/internal/
├── config/config.go           # 환경변수 로드
├── database/database.go       # GORM DB 연결
├── model/user.go              # User GORM 모델
├── repository/user.go         # User 레포지토리
├── service/auth.go            # 인증 서비스 (JWT, bcrypt, CSRF)
├── handler/auth.go            # 인증 핸들러
├── middleware/
│   ├── auth.go                # JWT 검증 미들웨어
│   └── csrf.go                # CSRF 검증 미들웨어
└── dto/
    ├── request.go             # 요청 DTO
    └── response.go            # 응답 DTO
```

### Frontend 구조

```
frontend/src/
├── lib/api/
│   ├── client.ts              # API 클라이언트
│   └── auth.ts                # 인증 API
├── stores/auth.ts             # Zustand 인증 스토어
├── components/auth/
│   └── auth-form.tsx          # 인증 폼 컴포넌트
├── app/
│   ├── (auth)/
│   │   ├── login/page.tsx     # 로그인 페이지
│   │   └── signup/page.tsx    # 회원가입 페이지
│   └── (protected)/
│       └── dashboard/page.tsx # 대시보드 페이지
└── middleware.ts              # Next.js 라우트 보호
```

### 추가된 패키지

**Backend**
- `gorm.io/gorm` - ORM
- `gorm.io/driver/postgres` - PostgreSQL 드라이버
- `github.com/golang-jwt/jwt/v5` - JWT
- `github.com/google/uuid` - UUID

**Frontend**
- `zustand` - 상태 관리

### API 엔드포인트

| Method | Endpoint | 설명 |
|--------|----------|------|
| POST | /api/v1/auth/signup | 회원가입 |
| POST | /api/v1/auth/login | 로그인 |
| POST | /api/v1/auth/logout | 로그아웃 |
| GET | /api/v1/auth/me | 현재 사용자 조회 (인증 필요) |
