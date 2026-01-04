package model

import (
	"time"

	"github.com/google/uuid"
)

type AnalysisStatus string
type AnalysisErrorCode string

const (
	AnalysisStatusPending    AnalysisStatus = "pending"
	AnalysisStatusProcessing AnalysisStatus = "processing"
	AnalysisStatusCompleted  AnalysisStatus = "completed"
	AnalysisStatusFailed     AnalysisStatus = "failed"
)

const (
	AnalysisErrorCodeMLModel   AnalysisErrorCode = "ML_MODEL_ERROR"
	AnalysisErrorCodeOpenAI    AnalysisErrorCode = "OPENAI_API_ERROR"
	AnalysisErrorCodeInvalidIn AnalysisErrorCode = "INVALID_INPUT"
	AnalysisErrorCodeTimeout   AnalysisErrorCode = "TIMEOUT"
	AnalysisErrorCodeInternal  AnalysisErrorCode = "INTERNAL_ERROR"
)

type Analysis struct {
	ID           uuid.UUID          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WritingID    uuid.UUID          `gorm:"type:uuid;not null;index" json:"writing_id"`
	TaskID       *uuid.UUID         `gorm:"type:uuid;uniqueIndex" json:"task_id"`
	Status       AnalysisStatus     `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	AIScore      *float64           `gorm:"type:decimal(5,2)" json:"ai_score"`
	Feedback     *string            `gorm:"type:text" json:"feedback"`
	ErrorCode    *AnalysisErrorCode `gorm:"type:varchar(50)" json:"error_code"`
	ErrorMessage *string            `gorm:"type:text" json:"error_message"`
	LatencyMs    *int               `gorm:"type:integer" json:"latency_ms"`
	RetryCount   int                `gorm:"not null;default:0" json:"retry_count"`
	CreatedAt    time.Time          `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time          `gorm:"not null;default:now()" json:"updated_at"`

	// Relations
	Writing Writing `gorm:"foreignKey:WritingID" json:"-"`
}

func (Analysis) TableName() string {
	return "analyses"
}

type AnalysisLog struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AnalysisID   uuid.UUID `gorm:"type:uuid;not null;index" json:"analysis_id"`
	InputText    string    `gorm:"type:text;not null" json:"input_text"`
	ModelVersion string    `gorm:"type:varchar(50);not null" json:"model_version"`
	RawOutput    *string   `gorm:"type:jsonb" json:"raw_output"`
	CreatedAt    time.Time `gorm:"not null;default:now()" json:"created_at"`

	// Relations
	Analysis Analysis `gorm:"foreignKey:AnalysisID" json:"-"`
}

func (AnalysisLog) TableName() string {
	return "analysis_logs"
}
