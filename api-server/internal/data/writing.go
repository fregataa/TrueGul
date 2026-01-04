package data

import (
	"time"

	"github.com/google/uuid"
)

type WritingType string
type WritingStatus string

const (
	WritingTypeEssay       WritingType = "essay"
	WritingTypeCoverLetter WritingType = "cover_letter"
)

const (
	WritingStatusDraft     WritingStatus = "draft"
	WritingStatusSubmitted WritingStatus = "submitted"
	WritingStatusAnalyzed  WritingStatus = "analyzed"
)

type Writing struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Type        WritingType
	Title       string
	Content     string
	Status      WritingStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SubmittedAt *time.Time
}
