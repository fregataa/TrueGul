package service

import (
	"github.com/google/uuid"
	"github.com/truegul/api-server/internal/data"
	apperrors "github.com/truegul/api-server/internal/errors"
	"github.com/truegul/api-server/internal/repository"
)

const MaxContentLength = 2000

type WritingService struct {
	writingRepo *repository.WritingRepository
}

func NewWritingService(writingRepo *repository.WritingRepository) *WritingService {
	return &WritingService{writingRepo: writingRepo}
}

func (s *WritingService) Create(userID uuid.UUID, writingType, title, content string) (*data.Writing, error) {
	if len([]rune(content)) > MaxContentLength {
		return nil, apperrors.ContentTooLong("Content exceeds maximum length of 2000 characters")
	}

	writing := &data.Writing{
		UserID:  userID,
		Type:    data.WritingType(writingType),
		Title:   title,
		Content: content,
		Status:  data.WritingStatusDraft,
	}

	if err := s.writingRepo.Create(writing); err != nil {
		return nil, err
	}

	return writing, nil
}

func (s *WritingService) GetByID(id, userID uuid.UUID) (*data.Writing, error) {
	writing, err := s.writingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if writing.UserID != userID {
		return nil, apperrors.Forbidden("Access denied")
	}

	return writing, nil
}

func (s *WritingService) List(userID uuid.UUID, page, limit int) ([]*data.Writing, int64, error) {
	offset := (page - 1) * limit
	return s.writingRepo.FindByUserID(userID, offset, limit)
}

func (s *WritingService) Update(id, userID uuid.UUID, writingType, title, content *string) (*data.Writing, error) {
	writing, err := s.writingRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if writing.UserID != userID {
		return nil, apperrors.Forbidden("Access denied")
	}

	if writingType != nil {
		writing.Type = data.WritingType(*writingType)
	}
	if title != nil {
		writing.Title = *title
	}
	if content != nil {
		if len([]rune(*content)) > MaxContentLength {
			return nil, apperrors.ContentTooLong("Content exceeds maximum length of 2000 characters")
		}
		writing.Content = *content
	}

	if err := s.writingRepo.Update(writing); err != nil {
		return nil, err
	}

	return writing, nil
}

func (s *WritingService) Delete(id, userID uuid.UUID) error {
	writing, err := s.writingRepo.FindByID(id)
	if err != nil {
		return err
	}

	if writing.UserID != userID {
		return apperrors.Forbidden("Access denied")
	}

	return s.writingRepo.Delete(id)
}
