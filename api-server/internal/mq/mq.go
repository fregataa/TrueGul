package mq

import (
	"context"

	"github.com/google/uuid"
)

type WritingType string

const (
	WritingTypeEssay       WritingType = "essay"
	WritingTypeCoverLetter WritingType = "cover_letter"
)

type AnalysisTask struct {
	Version     string      `json:"version"`
	TaskID      uuid.UUID   `json:"task_id"`
	WritingID   uuid.UUID   `json:"writing_id"`
	Content     string      `json:"content"`
	WritingType WritingType `json:"writing_type"`
	CallbackURL string      `json:"callback_url"`
}

type Publisher interface {
	Publish(ctx context.Context, task AnalysisTask) error
	Close() error
}
