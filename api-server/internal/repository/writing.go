package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/truegul/api-server/internal/data"
	apperrors "github.com/truegul/api-server/internal/errors"
	"github.com/truegul/api-server/internal/model"
	"gorm.io/gorm"
)

type WritingRepository struct {
	db *gorm.DB
}

func NewWritingRepository(db *gorm.DB) *WritingRepository {
	return &WritingRepository{db: db}
}

func (r *WritingRepository) Create(writing *data.Writing) error {
	m := toWritingModel(writing)
	if err := r.db.Create(m).Error; err != nil {
		return apperrors.InternalServerWrap(err, "Failed to create writing")
	}
	writing.ID = m.ID
	writing.CreatedAt = m.CreatedAt
	writing.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *WritingRepository) FindByID(id uuid.UUID) (*data.Writing, error) {
	var m model.Writing
	err := r.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("Writing not found")
		}
		return nil, apperrors.InternalServerWrap(err, "Failed to find writing")
	}
	return toWritingData(&m), nil
}

func (r *WritingRepository) FindByUserID(userID uuid.UUID, offset, limit int) ([]*data.Writing, int64, error) {
	var writings []model.Writing
	var total int64

	query := r.db.Model(&model.Writing{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, apperrors.InternalServerWrap(err, "Failed to count writings")
	}

	if err := query.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&writings).Error; err != nil {
		return nil, 0, apperrors.InternalServerWrap(err, "Failed to list writings")
	}

	result := make([]*data.Writing, len(writings))
	for i, w := range writings {
		result[i] = toWritingData(&w)
	}

	return result, total, nil
}

func (r *WritingRepository) Update(writing *data.Writing) error {
	m := toWritingModel(writing)
	if err := r.db.Save(m).Error; err != nil {
		return apperrors.InternalServerWrap(err, "Failed to update writing")
	}
	writing.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *WritingRepository) Delete(id uuid.UUID) error {
	if err := r.db.Delete(&model.Writing{}, "id = ?", id).Error; err != nil {
		return apperrors.InternalServerWrap(err, "Failed to delete writing")
	}
	return nil
}

func toWritingModel(d *data.Writing) *model.Writing {
	return &model.Writing{
		ID:          d.ID,
		UserID:      d.UserID,
		Type:        model.WritingType(d.Type),
		Title:       d.Title,
		Content:     d.Content,
		Status:      model.WritingStatus(d.Status),
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		SubmittedAt: d.SubmittedAt,
	}
}

func toWritingData(m *model.Writing) *data.Writing {
	return &data.Writing{
		ID:          m.ID,
		UserID:      m.UserID,
		Type:        data.WritingType(m.Type),
		Title:       m.Title,
		Content:     m.Content,
		Status:      data.WritingStatus(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		SubmittedAt: m.SubmittedAt,
	}
}
