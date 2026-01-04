package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	apperrors "github.com/truegul/api-server/internal/errors"
	"github.com/truegul/api-server/internal/model"
)

type AnalysisRepository struct {
	db *gorm.DB
}

func NewAnalysisRepository(db *gorm.DB) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

func (r *AnalysisRepository) Create(analysis *model.Analysis) error {
	if err := r.db.Create(analysis).Error; err != nil {
		return apperrors.InternalServerWrap(err, "Failed to create analysis")
	}
	return nil
}

func (r *AnalysisRepository) FindByID(id uuid.UUID) (*model.Analysis, error) {
	var analysis model.Analysis
	if err := r.db.First(&analysis, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound("Analysis not found")
		}
		return nil, apperrors.InternalServerWrap(err, "Failed to find analysis")
	}
	return &analysis, nil
}

func (r *AnalysisRepository) FindByTaskID(taskID uuid.UUID) (*model.Analysis, error) {
	var analysis model.Analysis
	if err := r.db.First(&analysis, "task_id = ?", taskID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound("Analysis not found")
		}
		return nil, apperrors.InternalServerWrap(err, "Failed to find analysis")
	}
	return &analysis, nil
}

func (r *AnalysisRepository) FindByWritingID(writingID uuid.UUID) (*model.Analysis, error) {
	var analysis model.Analysis
	if err := r.db.Where("writing_id = ?", writingID).Order("created_at DESC").First(&analysis).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.NotFound("Analysis not found")
		}
		return nil, apperrors.InternalServerWrap(err, "Failed to find analysis")
	}
	return &analysis, nil
}

func (r *AnalysisRepository) Update(analysis *model.Analysis) error {
	if err := r.db.Save(analysis).Error; err != nil {
		return apperrors.InternalServerWrap(err, "Failed to update analysis")
	}
	return nil
}

func (r *AnalysisRepository) UpdateResult(taskID uuid.UUID, status model.AnalysisStatus, aiScore *float64, feedback *string, errorCode *model.AnalysisErrorCode, errorMessage *string, latencyMs *int) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if aiScore != nil {
		updates["ai_score"] = *aiScore
	}
	if feedback != nil {
		updates["feedback"] = *feedback
	}
	if errorCode != nil {
		updates["error_code"] = *errorCode
	}
	if errorMessage != nil {
		updates["error_message"] = *errorMessage
	}
	if latencyMs != nil {
		updates["latency_ms"] = *latencyMs
	}

	result := r.db.Model(&model.Analysis{}).Where("task_id = ?", taskID).Updates(updates)
	if result.Error != nil {
		return apperrors.InternalServerWrap(result.Error, "Failed to update analysis result")
	}
	if result.RowsAffected == 0 {
		return apperrors.NotFound("Analysis not found")
	}
	return nil
}

func (r *AnalysisRepository) IncrementRetryCount(taskID uuid.UUID) error {
	result := r.db.Model(&model.Analysis{}).Where("task_id = ?", taskID).Update("retry_count", gorm.Expr("retry_count + 1"))
	if result.Error != nil {
		return apperrors.InternalServerWrap(result.Error, "Failed to increment retry count")
	}
	return nil
}

func (r *AnalysisRepository) CreateLog(log *model.AnalysisLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return apperrors.InternalServerWrap(err, "Failed to create analysis log")
	}
	return nil
}
