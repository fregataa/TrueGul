# Sprint 1: Backend - API Server

## Overview

| 항목 | 내용 |
|------|------|
| 목표 | 채점 요청/결과 API 구현 |
| 선행 조건 | Sprint 0 완료 |
| 병렬 진행 | Sprint 2 (ML Server) |
| 후속 Sprint | Sprint 4 (Integration) |

---

## Tasks

| ID | Task | 파일 | 상태 | 비고 |
|----|------|------|------|------|
| S1-1 | DB 마이그레이션 작성 | `migrations/00X_*.sql` | TODO | |
| S1-2 | Submission 모델/리포지토리 | `internal/repository/submission.go` | TODO | |
| S1-3 | ScoringResult 모델/리포지토리 | `internal/repository/scoring_result.go` | TODO | |
| S1-4 | PushToken 모델/리포지토리 | `internal/repository/push_token.go` | TODO | |
| S1-5 | Submission 서비스 | `internal/service/submission.go` | TODO | Redis 큐 연동 |
| S1-6 | Push 서비스 | `internal/service/push.go` | TODO | FCM/APNs 통합 |
| S1-7 | Submission 핸들러 | `internal/handler/submission.go` | TODO | |
| S1-8 | Push 핸들러 | `internal/handler/push.go` | TODO | |
| S1-9 | 라우터 등록 | `internal/router/router.go` | TODO | |
| S1-10 | Callback 엔드포인트 수정 | `internal/handler/callback.go` | TODO | 새 결과 형식 + Push 발송 |

---

## S1-1: Database Schema

### submissions 테이블

```sql
CREATE TABLE submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    question_type VARCHAR(10) NOT NULL DEFAULT '54',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_status ON submissions(status);
CREATE INDEX idx_submissions_created_at ON submissions(created_at DESC);

COMMENT ON COLUMN submissions.status IS 'pending | processing | completed | failed';
COMMENT ON COLUMN submissions.question_type IS 'TOPIK 문항 번호 (54, 53, 52, 51)';
```

### scoring_results 테이블

```sql
CREATE TABLE scoring_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,

    -- 채점 결과
    content_score INTEGER CHECK (content_score >= 0 AND content_score <= 20),
    structure_score INTEGER CHECK (structure_score >= 0 AND structure_score <= 15),
    language_score INTEGER CHECK (language_score >= 0 AND language_score <= 15),
    total_score INTEGER CHECK (total_score >= 0 AND total_score <= 50),

    -- 피드백
    content_feedback TEXT,
    structure_feedback TEXT,
    language_feedback TEXT,
    overall_feedback TEXT,

    -- 추정 레벨
    level_estimate VARCHAR(10),

    -- AI 감지
    ai_detection_score DECIMAL(5,4),
    ai_detection_flagged BOOLEAN DEFAULT FALSE,
    ai_detection_model VARCHAR(50),

    -- 메타
    llm_model VARCHAR(50),
    llm_tokens_used INTEGER,
    processing_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(submission_id)
);

CREATE INDEX idx_scoring_results_submission_id ON scoring_results(submission_id);
```

### push_tokens 테이블

```sql
CREATE TABLE push_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    platform VARCHAR(10) NOT NULL CHECK (platform IN ('ios', 'android')),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(user_id, platform)
);

CREATE INDEX idx_push_tokens_user_id ON push_tokens(user_id);
CREATE INDEX idx_push_tokens_active ON push_tokens(is_active) WHERE is_active = TRUE;
```

---

## S1-2 ~ S1-4: Models & Repositories

### Submission Model

```go
// internal/model/submission.go
type Submission struct {
    ID           uuid.UUID  `json:"id"`
    UserID       uuid.UUID  `json:"user_id"`
    Content      string     `json:"content"`
    QuestionType string     `json:"question_type"`
    Status       string     `json:"status"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

type SubmissionStatus string

const (
    StatusPending    SubmissionStatus = "pending"
    StatusProcessing SubmissionStatus = "processing"
    StatusCompleted  SubmissionStatus = "completed"
    StatusFailed     SubmissionStatus = "failed"
)
```

### ScoringResult Model

```go
// internal/model/scoring_result.go
type ScoringResult struct {
    ID                  uuid.UUID `json:"id"`
    SubmissionID        uuid.UUID `json:"submission_id"`

    // Scores
    ContentScore        int       `json:"content_score"`
    StructureScore      int       `json:"structure_score"`
    LanguageScore       int       `json:"language_score"`
    TotalScore          int       `json:"total_score"`

    // Feedback
    ContentFeedback     string    `json:"content_feedback"`
    StructureFeedback   string    `json:"structure_feedback"`
    LanguageFeedback    string    `json:"language_feedback"`
    OverallFeedback     string    `json:"overall_feedback"`

    LevelEstimate       string    `json:"level_estimate"`

    // AI Detection
    AIDetectionScore    float64   `json:"ai_detection_score"`
    AIDetectionFlagged  bool      `json:"ai_detection_flagged"`

    CreatedAt           time.Time `json:"created_at"`
}
```

### Repository Interface

```go
// internal/repository/submission_repository.go
type SubmissionRepository interface {
    Create(ctx context.Context, s *model.Submission) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.Submission, error)
    GetByUserID(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]model.Submission, int, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

// internal/repository/scoring_result_repository.go
type ScoringResultRepository interface {
    Create(ctx context.Context, r *model.ScoringResult) error
    GetBySubmissionID(ctx context.Context, submissionID uuid.UUID) (*model.ScoringResult, error)
}

// internal/repository/push_token_repository.go
type PushTokenRepository interface {
    Upsert(ctx context.Context, userID uuid.UUID, token, platform string) error
    GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.PushToken, error)
    Delete(ctx context.Context, userID uuid.UUID, platform string) error
}
```

---

## S1-5: Submission Service

```go
// internal/service/submission_service.go
type SubmissionService struct {
    repo        repository.SubmissionRepository
    resultRepo  repository.ScoringResultRepository
    redisClient *redis.Client
    pushService *PushService
}

func (s *SubmissionService) Submit(ctx context.Context, userID uuid.UUID, content string) (*model.Submission, error) {
    // 1. Validate content length (600-700자 권장, 최대 800자)
    if len([]rune(content)) > 800 {
        return nil, apperrors.BadRequest("Content too long")
    }

    // 2. Create submission
    submission := &model.Submission{
        UserID:       userID,
        Content:      content,
        QuestionType: "54",
        Status:       string(model.StatusPending),
    }
    if err := s.repo.Create(ctx, submission); err != nil {
        return nil, err
    }

    // 3. Enqueue to Redis
    task := map[string]interface{}{
        "submission_id": submission.ID,
        "content":       content,
        "user_id":       userID,
    }
    if err := s.redisClient.LPush(ctx, "topik_scoring_queue", task).Err(); err != nil {
        return nil, err
    }

    // 4. Update status to processing
    s.repo.UpdateStatus(ctx, submission.ID, string(model.StatusProcessing))

    return submission, nil
}

func (s *SubmissionService) GetResult(ctx context.Context, userID, submissionID uuid.UUID) (*dto.SubmissionWithResult, error) {
    // Verify ownership and return result
}

func (s *SubmissionService) GetHistory(ctx context.Context, userID uuid.UUID, status string, page, limit int) (*dto.SubmissionList, error) {
    // Return paginated history
}
```

---

## S1-6: Push Service

```go
// internal/service/push_service.go
type PushService struct {
    repo      repository.PushTokenRepository
    fcmClient *messaging.Client  // Firebase Admin SDK
    apnsClient *apns2.Client     // APNs client
}

func (s *PushService) RegisterToken(ctx context.Context, userID uuid.UUID, token, platform string) error {
    return s.repo.Upsert(ctx, userID, token, platform)
}

func (s *PushService) SendScoringComplete(ctx context.Context, userID, submissionID uuid.UUID, totalScore int) error {
    tokens, err := s.repo.GetByUserID(ctx, userID)
    if err != nil {
        return err
    }

    for _, t := range tokens {
        notification := &Notification{
            Title: "채점 완료",
            Body:  fmt.Sprintf("총점 %d점을 받았습니다. 상세 피드백을 확인하세요.", totalScore),
            Data: map[string]string{
                "type":          "scoring_complete",
                "submission_id": submissionID.String(),
            },
        }

        switch t.Platform {
        case "android":
            s.sendFCM(t.Token, notification)
        case "ios":
            s.sendAPNs(t.Token, notification)
        }
    }
    return nil
}
```

---

## S1-7 ~ S1-8: Handlers

### API Endpoints

| Method | Endpoint | Handler | Auth |
|--------|----------|---------|------|
| POST | `/api/v1/submissions` | CreateSubmission | Required |
| GET | `/api/v1/submissions/:id` | GetSubmission | Required |
| GET | `/api/v1/submissions` | ListSubmissions | Required |
| POST | `/api/v1/push/register` | RegisterPushToken | Required |
| DELETE | `/api/v1/push/unregister` | UnregisterPushToken | Required |

### Request/Response DTOs

```go
// internal/dto/submission.go

// CreateSubmissionRequest
type CreateSubmissionRequest struct {
    Content      string `json:"content" binding:"required,min=100,max=800"`
    QuestionType string `json:"question_type,omitempty"`
}

// SubmissionResponse
type SubmissionResponse struct {
    ID           uuid.UUID  `json:"id"`
    Status       string     `json:"status"`
    QuestionType string     `json:"question_type"`
    CreatedAt    time.Time  `json:"created_at"`
}

// SubmissionWithResultResponse
type SubmissionWithResultResponse struct {
    Submission SubmissionResponse `json:"submission"`
    Result     *ScoringResultResponse `json:"result,omitempty"`
}

// ScoringResultResponse
type ScoringResultResponse struct {
    Scores   ScoresResponse   `json:"scores"`
    Feedback FeedbackResponse `json:"feedback"`
    Level    string           `json:"level_estimate"`
    AIDetection AIDetectionResponse `json:"ai_detection"`
}

// RegisterPushTokenRequest
type RegisterPushTokenRequest struct {
    Token    string `json:"token" binding:"required"`
    Platform string `json:"platform" binding:"required,oneof=ios android"`
}
```

---

## S1-10: Callback Handler Update

```go
// internal/handler/callback.go
func (h *CallbackHandler) HandleScoringResult(c *gin.Context) {
    var req dto.ScoringCallbackRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        handleValidationError(c, err.Error())
        return
    }

    // 1. Verify callback secret
    if !h.verifySecret(c.GetHeader("X-Callback-Secret")) {
        handleError(c, apperrors.Unauthorized("Invalid callback secret"))
        return
    }

    // 2. Save scoring result
    result := &model.ScoringResult{
        SubmissionID:       req.SubmissionID,
        ContentScore:       req.Scores.Content,
        StructureScore:     req.Scores.Structure,
        LanguageScore:      req.Scores.Language,
        TotalScore:         req.Scores.Total,
        ContentFeedback:    req.Feedback.Content,
        StructureFeedback:  req.Feedback.Structure,
        LanguageFeedback:   req.Feedback.Language,
        OverallFeedback:    req.Feedback.Overall,
        LevelEstimate:      req.LevelEstimate,
        AIDetectionScore:   req.AIDetection.Score,
        AIDetectionFlagged: req.AIDetection.Flagged,
    }
    if err := h.resultRepo.Create(c, result); err != nil {
        handleError(c, err)
        return
    }

    // 3. Update submission status
    h.submissionRepo.UpdateStatus(c, req.SubmissionID, "completed")

    // 4. Send push notification
    submission, _ := h.submissionRepo.GetByID(c, req.SubmissionID)
    h.pushService.SendScoringComplete(c, submission.UserID, req.SubmissionID, req.Scores.Total)

    c.JSON(http.StatusOK, dto.MessageResponse{Message: "Result saved"})
}
```

---

## Completion Criteria

- [ ] DB 마이그레이션 작성 및 적용
- [ ] 모든 모델/리포지토리 구현
- [ ] Submission 서비스 구현 및 Redis 연동
- [ ] Push 서비스 구현 (FCM/APNs)
- [ ] 모든 API 엔드포인트 구현
- [ ] Callback 핸들러 수정
- [ ] 단위 테스트 작성
- [ ] API 문서 업데이트

---

*Sprint 2 (ML Server)와 병렬 진행 가능*
