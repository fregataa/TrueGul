# Version 1: OCR 도입

## Overview

| 항목 | 내용 |
|------|------|
| 목표 | 수기 답안 이미지 업로드 및 OCR 변환 지원 |
| 선행 조건 | v0 안정화 완료 |
| 핵심 가치 | 실전 TOPIK 환경 시뮬레이션 |

---

## Background

v0에서는 텍스트 직접 입력만 지원하지만, 실제 TOPIK 시험은 수기로 작성합니다.
v1에서는 손글씨 답안 사진을 업로드하여 OCR로 텍스트 변환 후 채점하는 기능을 추가합니다.

---

## Features

### 1. 이미지 업로드

| 기능 | 상세 |
|------|------|
| 촬영 | 앱 내 카메라로 답안지 촬영 |
| 갤러리 | 기존 사진 선택 |
| 이미지 편집 | 회전, 크롭, 밝기 조절 |
| 미리보기 | 업로드 전 확인 |

### 2. OCR 처리

| 항목 | 내용 |
|------|------|
| 기술 | 외부 API (Google Vision / Naver Clova OCR) |
| 대상 | 한글 손글씨 |
| 처리 시간 | 예상 5-10초 |

### 3. OCR 결과 검토/수정

```
┌─────────────────────────────────────┐
│ ←         OCR 결과 확인             │
├─────────────────────────────────────┤
│                                     │
│  [원본 이미지 썸네일]               │
│                                     │
├─────────────────────────────────────┤
│                                     │
│  인식된 텍스트:                     │
│  ┌─────────────────────────────┐   │
│  │                             │   │
│  │  [편집 가능한 텍스트 영역]  │   │
│  │                             │   │
│  │  환경 보호는 현대 사회에서  │   │
│  │  매우 중요한 문제입니다...  │   │
│  │                             │   │
│  └─────────────────────────────┘   │
│                                     │
│  ⚠️ 인식이 정확하지 않을 수 있습니다│
│     직접 수정해주세요.              │
│                                     │
│  글자 수: 623 / 700                 │
│                                     │
├─────────────────────────────────────┤
│                                     │
│  [        채점 요청        ]        │
│                                     │
└─────────────────────────────────────┘
```

| 기능 | 상세 |
|------|------|
| 텍스트 편집 | OCR 오류 직접 수정 |
| 원본 비교 | 이미지와 텍스트 동시 확인 |
| 글자 수 표시 | 실시간 카운터 |

---

## Technical Implementation

### API Endpoints

| Method | Endpoint | 설명 |
|--------|----------|------|
| POST | `/api/v1/submissions/image` | 이미지 업로드 + OCR |
| GET | `/api/v1/submissions/:id/image` | 원본 이미지 조회 |

### Request Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  Mobile  │────▶│   API    │────▶│    S3    │     │   OCR    │
│   App    │     │  Server  │     │  Upload  │     │   API    │
└──────────┘     └────┬─────┘     └──────────┘     └────┬─────┘
                      │                                  │
                      │◀─────────────────────────────────┘
                      │         OCR 결과 반환
                      ▼
               ┌──────────┐
               │  Mobile  │  텍스트 편집 화면
               │   App    │
               └────┬─────┘
                    │
                    ▼  수정 후 제출
               (기존 채점 플로우)
```

### Database Schema

```sql
-- submissions 테이블에 컬럼 추가
ALTER TABLE submissions ADD COLUMN input_type VARCHAR(10) DEFAULT 'text';
ALTER TABLE submissions ADD COLUMN image_url TEXT;
ALTER TABLE submissions ADD COLUMN ocr_raw_text TEXT;
ALTER TABLE submissions ADD COLUMN ocr_confidence DECIMAL(5,4);

COMMENT ON COLUMN submissions.input_type IS 'text | image';
COMMENT ON COLUMN submissions.ocr_raw_text IS 'OCR 원본 결과 (수정 전)';
COMMENT ON COLUMN submissions.ocr_confidence IS 'OCR 신뢰도 점수';
```

### OCR API Integration

```python
# app/services/ocr_service.py
from abc import ABC, abstractmethod

class OCRService(ABC):
    @abstractmethod
    async def recognize(self, image_bytes: bytes) -> OCRResult:
        pass

class GoogleVisionOCR(OCRService):
    async def recognize(self, image_bytes: bytes) -> OCRResult:
        # Google Cloud Vision API 호출
        pass

class NaverClovaOCR(OCRService):
    async def recognize(self, image_bytes: bytes) -> OCRResult:
        # Naver Clova OCR API 호출
        pass

class OCRResult:
    text: str
    confidence: float
    word_boxes: list[WordBox]  # 단어별 위치 정보
```

---

## OCR API Comparison

| 항목 | Google Vision | Naver Clova OCR |
|------|--------------|-----------------|
| 한글 손글씨 | 보통 | 우수 |
| 가격 | $1.50/1000건 | 월 300건 무료 |
| 응답 속도 | 빠름 | 빠름 |
| API 안정성 | 높음 | 높음 |

**권장**: Naver Clova OCR (한글 특화)

---

## User Experience

### 플로우

```
1. Submit 화면에서 "이미지로 제출" 선택
2. 카메라 촬영 또는 갤러리 선택
3. 이미지 크롭/회전 (필요시)
4. 업로드 → OCR 처리 (로딩)
5. OCR 결과 확인 화면
6. 텍스트 수정 (필요시)
7. 채점 요청
8. (이후 기존 플로우와 동일)
```

### 에러 처리

| 상황 | 대응 |
|------|------|
| OCR 실패 | "인식 실패. 다시 촬영해주세요" |
| 낮은 신뢰도 | "인식률이 낮습니다. 확인 후 수정해주세요" |
| 이미지 품질 불량 | "이미지가 흐립니다. 다시 촬영해주세요" |
| 텍스트 없음 | "텍스트를 찾을 수 없습니다" |

---

## Cost Estimation

| 항목 | 월 예상 (MAU 1,000) |
|------|---------------------|
| OCR API (10% 이미지 사용) | $3-5 |
| S3 스토리지 (이미지) | $1-2 |
| S3 전송 비용 | $1 |
| **Total** | **~$5-8** |

---

## Success Metrics

| 지표 | 목표 |
|------|------|
| OCR 인식률 | 90% 이상 (수정 없이 사용 가능) |
| 이미지 제출 비율 | 전체 제출의 20% 이상 |
| OCR 처리 시간 | 10초 이내 |
| 사용자 만족도 | 4.0/5.0 이상 |

---

## Risks & Mitigation

| 리스크 | 대응 |
|--------|------|
| OCR 인식률 저조 | 사용자 수정 UI 제공, 가이드 개선 |
| 이미지 용량 과다 | 클라이언트에서 압축/리사이즈 |
| 촬영 환경 불량 | 촬영 가이드 제공 (조명, 각도) |

---

## Implementation Phases

### Phase 1: 기본 OCR
- 이미지 업로드
- OCR API 연동
- 결과 확인/수정 UI

### Phase 2: 개선
- 이미지 전처리 (밝기, 대비 자동 조절)
- 촬영 가이드 오버레이
- 단어별 신뢰도 표시

---

## Dependencies

- v0 완료 및 안정화
- OCR API 계정 및 결제 설정
- S3 이미지 저장 설정
- 모바일 카메라/갤러리 권한

---

*v1 완료 후 v2 (데이터 축적) 진행*
