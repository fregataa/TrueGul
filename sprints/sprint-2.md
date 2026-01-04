# Sprint 2: 글 작성 기능

## 목표

글 CRUD 및 임시 저장 기능 구현

---

## 작업 목록

### Backend (API Server)

- [x] 글 모델 및 레포지토리 구현
- [x] GET /api/writings 구현 (목록 조회, 페이지네이션)
- [x] GET /api/writings/:id 구현 (상세 조회)
- [x] POST /api/writings 구현 (글 생성/임시 저장)
- [x] PUT /api/writings/:id 구현 (글 수정)
- [x] DELETE /api/writings/:id 구현 (글 삭제)
- [x] 글 길이 검증 (최대 2,000자)
- [x] 글 종류 검증 (essay, cover_letter)
- [x] 본인 글만 접근 가능하도록 권한 체크

### Frontend

- [x] 글 목록 페이지 (대시보드)
- [x] 글 작성/수정 페이지
- [x] 텍스트 에디터 컴포넌트 (글자 수 표시)
- [x] 글 종류 선택 UI
- [x] 임시 저장 기능 (자동 저장 또는 수동)
- [x] 글 삭제 확인 모달

### 테스트

- [ ] 글 CRUD API 단위 테스트
- [ ] 권한 체크 테스트

---

## 완료 조건

- [x] 글 생성, 조회, 수정, 삭제 가능
- [x] 임시 저장 기능 동작
- [x] 2,000자 초과 시 에러 반환
- [x] 다른 사용자의 글 접근 시 403 반환
- [ ] 글 CRUD 테스트 통과

---

## 구현 상세

### Backend 구조

```
api-server/internal/
├── model/writing.go           # Writing GORM 모델
├── data/writing.go            # Writing 도메인 구조체
├── repository/writing.go      # Writing 레포지토리
├── service/writing.go         # Writing 비즈니스 로직
├── handler/writing.go         # Writing HTTP 핸들러
└── dto/
    ├── request.go             # CreateWritingRequest, UpdateWritingRequest, ListWritingsQuery
    └── response.go            # WritingResponse, WritingListResponse
```

### Frontend 구조

```
frontend/src/
├── lib/api/writings.ts        # Writing API 클라이언트
├── stores/writings.ts         # Zustand Writing 스토어
├── components/writings/
│   └── writing-editor.tsx     # 글 작성/수정 에디터 컴포넌트
└── app/(protected)/
    ├── dashboard/page.tsx     # 글 목록 페이지
    └── writings/
        ├── new/page.tsx       # 새 글 작성 페이지
        └── [id]/
            ├── page.tsx       # 글 상세 보기 페이지
            └── edit/page.tsx  # 글 수정 페이지
```

### API 엔드포인트

| Method | Endpoint | 설명 |
|--------|----------|------|
| GET | /api/v1/writings | 글 목록 조회 (페이지네이션) |
| GET | /api/v1/writings/:id | 글 상세 조회 |
| POST | /api/v1/writings | 글 생성 |
| PUT | /api/v1/writings/:id | 글 수정 |
| DELETE | /api/v1/writings/:id | 글 삭제 |

### 주요 기능

1. **글 CRUD**: 글 생성, 조회, 수정, 삭제 기능
2. **글 종류**: 에세이(essay), 자기소개서(cover_letter)
3. **글 상태**: 임시저장(draft), 제출됨(submitted), 분석완료(analyzed)
4. **글자 수 제한**: 최대 2,000자 (한글 기준 유니코드 길이)
5. **자동 저장**: 수정 모드에서 30초마다 자동 저장
6. **권한 체크**: 본인 글만 접근 가능 (403 Forbidden)
7. **페이지네이션**: 글 목록 페이지네이션 지원
