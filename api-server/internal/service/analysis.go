package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/truegul/api-server/internal/config"
	"github.com/truegul/api-server/internal/data"
	apperrors "github.com/truegul/api-server/internal/errors"
	"github.com/truegul/api-server/internal/model"
	"github.com/truegul/api-server/internal/mq"
	"github.com/truegul/api-server/internal/repository"
)

const (
	MaxDailySubmissions = 5
	MaxRetries          = 3
)

type AnalysisService struct {
	analysisRepo *repository.AnalysisRepository
	writingRepo  *repository.WritingRepository
	userRepo     *repository.UserRepository
	publisher    mq.Publisher
	config       *config.Config
}

func NewAnalysisService(
	analysisRepo *repository.AnalysisRepository,
	writingRepo *repository.WritingRepository,
	userRepo *repository.UserRepository,
	publisher mq.Publisher,
	cfg *config.Config,
) *AnalysisService {
	return &AnalysisService{
		analysisRepo: analysisRepo,
		writingRepo:  writingRepo,
		userRepo:     userRepo,
		publisher:    publisher,
		config:       cfg,
	}
}

func (s *AnalysisService) SubmitWriting(ctx context.Context, writingID, userID uuid.UUID) (*model.Analysis, error) {
	writing, err := s.writingRepo.FindByID(writingID)
	if err != nil {
		return nil, err
	}

	if writing.UserID != userID {
		return nil, apperrors.Forbidden("You don't have permission to submit this writing")
	}

	if writing.Status != data.WritingStatusDraft {
		return nil, apperrors.Validation("Writing has already been submitted")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	if user.LastSubmitDate == nil || user.LastSubmitDate.Truncate(24*time.Hour).Before(today) {
		user.DailySubmitCount = 0
	}

	if user.DailySubmitCount >= MaxDailySubmissions {
		return nil, apperrors.New(
			apperrors.CodeForbidden,
			fmt.Sprintf("Daily submission limit reached (%d/%d)", user.DailySubmitCount, MaxDailySubmissions),
			429,
		)
	}

	taskID := uuid.New()
	analysis := &model.Analysis{
		WritingID: writingID,
		TaskID:    &taskID,
		Status:    model.AnalysisStatusPending,
	}

	if err := s.analysisRepo.Create(analysis); err != nil {
		return nil, err
	}

	now := time.Now()
	writing.Status = data.WritingStatusSubmitted
	writing.SubmittedAt = &now
	if err := s.writingRepo.Update(writing); err != nil {
		return nil, err
	}

	user.DailySubmitCount++
	user.LastSubmitDate = &now
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	callbackURL := fmt.Sprintf("%s%s", s.config.CallbackBaseURL, s.config.CallbackPath)
	task := mq.AnalysisTask{
		Version:     "1",
		TaskID:      taskID,
		WritingID:   writingID,
		Content:     writing.Content,
		WritingType: mq.WritingType(writing.Type),
		CallbackURL: callbackURL,
	}

	if err := s.publisher.Publish(ctx, task); err != nil {
		return nil, apperrors.InternalServerWrap(err, "Failed to queue analysis task")
	}

	return analysis, nil
}

func (s *AnalysisService) GetAnalysis(writingID, userID uuid.UUID) (*model.Analysis, error) {
	writing, err := s.writingRepo.FindByID(writingID)
	if err != nil {
		return nil, err
	}

	if writing.UserID != userID {
		return nil, apperrors.Forbidden("You don't have permission to view this analysis")
	}

	return s.analysisRepo.FindByWritingID(writingID)
}

func (s *AnalysisService) HandleCallback(taskID uuid.UUID, status string, result *CallbackResult, callbackErr *CallbackError) error {
	analysis, err := s.analysisRepo.FindByTaskID(taskID)
	if err != nil {
		return err
	}

	if analysis.Status == model.AnalysisStatusCompleted || analysis.Status == model.AnalysisStatusFailed {
		return nil
	}

	if status == "completed" && result != nil {
		aiScore := result.AIProbability
		feedback := result.Feedback
		latencyMs := result.LatencyMs

		if err := s.analysisRepo.UpdateResult(
			taskID,
			model.AnalysisStatusCompleted,
			&aiScore,
			&feedback,
			nil,
			nil,
			&latencyMs,
		); err != nil {
			return err
		}

		writing, err := s.writingRepo.FindByID(analysis.WritingID)
		if err != nil {
			return err
		}
		writing.Status = data.WritingStatusAnalyzed
		return s.writingRepo.Update(writing)
	}

	if status == "failed" && callbackErr != nil {
		errorCode := model.AnalysisErrorCode(callbackErr.Code)
		errorMessage := callbackErr.Message

		if callbackErr.Retryable && analysis.RetryCount < MaxRetries {
			if err := s.analysisRepo.IncrementRetryCount(taskID); err != nil {
				return err
			}
			return nil
		}

		return s.analysisRepo.UpdateResult(
			taskID,
			model.AnalysisStatusFailed,
			nil,
			nil,
			&errorCode,
			&errorMessage,
			nil,
		)
	}

	return apperrors.Validation("Invalid callback status")
}

type CallbackResult struct {
	AIProbability float64
	Feedback      string
	LatencyMs     int
}

type CallbackError struct {
	Code      string
	Message   string
	Retryable bool
}
