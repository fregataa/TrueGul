package model

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
	ID          uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	Type        WritingType   `gorm:"type:varchar(50);not null" json:"type"`
	Title       string        `gorm:"type:varchar(255);not null" json:"title"`
	Content     string        `gorm:"type:text;not null" json:"content"`
	Status      WritingStatus `gorm:"type:varchar(50);not null;default:'draft'" json:"status"`
	CreatedAt   time.Time     `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"not null;default:now()" json:"updated_at"`
	SubmittedAt *time.Time    `json:"submitted_at"`
}

func (Writing) TableName() string {
	return "writings"
}
