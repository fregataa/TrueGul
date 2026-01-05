# Sprint 4: Frontend 분석 UI 및 통합

## 목표

분석 결과 UI 및 전체 플로우 통합

---

## 작업 목록

### Frontend

- [x] 글 제출 버튼 및 확인 모달
- [x] 분석 진행 중 UI (로딩 상태, Polling)
- [x] 분석 결과 페이지
  - [x] AI 탐지 점수 시각화 (숫자 + 색상 배지)
  - [x] 품질 피드백 표시
- [x] Rate Limit 초과 시 안내 UI
- [x] 분석 실패 시 에러 UI
- [x] 제출 이력 표시 (글 목록에서 상태 확인)

### 통합 테스트 (별도 진행)

- [ ] 전체 플로우 E2E 테스트 (회원가입 → 글 작성 → 제출 → 결과 확인)
- [ ] 에러 케이스 테스트

### 문서화 (별도 진행)

- [ ] API 문서 정리 (OpenAPI/Swagger)
- [ ] 로컬 개발 환경 README 작성

---

## 완료 조건

- [x] 글 제출 후 분석 결과 확인 가능
- [x] 분석 진행 중 로딩 UI 표시
- [x] 분석 실패 시 적절한 에러 메시지 표시
- [x] Rate Limit 초과 시 안내 메시지 표시
- [ ] E2E 테스트 통과 (별도 진행)
- [ ] API 문서 완성 (별도 진행)

---

## 구현 내역

### 신규 파일

| 파일 | 설명 |
|------|------|
| `src/lib/api/analysis.ts` | Analysis API 클라이언트 |
| `src/stores/analysis.ts` | Analysis Zustand 스토어 (Polling 포함) |
| `src/components/ui/dialog.tsx` | shadcn Dialog 컴포넌트 |
| `src/components/ui/alert.tsx` | shadcn Alert 컴포넌트 |
| `src/components/ui/card.tsx` | shadcn Card 컴포넌트 |
| `src/components/ui/badge.tsx` | shadcn Badge 컴포넌트 |
| `src/components/analysis/submit-confirm-dialog.tsx` | 제출 확인 모달 |
| `src/components/analysis/analysis-progress.tsx` | 분석 진행 중 UI |
| `src/components/analysis/ai-score-badge.tsx` | AI 점수 배지 |
| `src/components/analysis/analysis-result.tsx` | 분석 결과 카드 |
| `src/components/analysis/rate-limit-warning.tsx` | Rate Limit 경고 |
| `src/components/analysis/analysis-error.tsx` | 분석 에러 UI |

### 수정된 파일

| 파일 | 변경 내용 |
|------|----------|
| `src/app/(protected)/writings/[id]/page.tsx` | 제출 버튼, 분석 결과 UI 통합 |
| `src/app/(protected)/dashboard/page.tsx` | Badge 컴포넌트로 상태 표시 개선, "분석 중" 상태 추가 |

### 주요 기능

1. **제출 확인 모달**: 하루 5회 제한 안내, 제출 후 수정 불가 안내
2. **Polling**: 지수 백오프 (2초 → 10초), 최대 60회 시도
3. **AI 점수 배지**: 0-30% 녹색(안전), 31-60% 노랑(주의), 61-100% 빨강(위험)
4. **에러 처리**: Rate Limit, ML 에러, OpenAI 에러, 타임아웃 등 코드별 메시지
