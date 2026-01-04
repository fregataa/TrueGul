package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/truegul/api-server/internal/config"
	"github.com/truegul/api-server/internal/dto"
	"github.com/truegul/api-server/internal/model"
	"github.com/truegul/api-server/internal/service"
)

const CallbackSecretHeader = "X-Callback-Secret"

type AnalysisHandler struct {
	analysisService *service.AnalysisService
	config          *config.Config
}

func NewAnalysisHandler(analysisService *service.AnalysisService, cfg *config.Config) *AnalysisHandler {
	return &AnalysisHandler{
		analysisService: analysisService,
		config:          cfg,
	}
}

func (h *AnalysisHandler) Submit(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			ErrorCode: "UNAUTHORIZED",
			Message:   "User not authenticated",
		})
		return
	}

	writingIDStr := c.Param("id")
	writingID, err := uuid.Parse(writingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			ErrorCode: "VALIDATION_ERROR",
			Message:   "Invalid writing ID",
		})
		return
	}

	analysis, err := h.analysisService.SubmitWriting(c.Request.Context(), writingID, userID.(uuid.UUID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, dto.SubmitResponse{
		Message:    "Analysis task queued",
		AnalysisID: analysis.ID,
	})
}

func (h *AnalysisHandler) GetAnalysis(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			ErrorCode: "UNAUTHORIZED",
			Message:   "User not authenticated",
		})
		return
	}

	writingIDStr := c.Param("id")
	writingID, err := uuid.Parse(writingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			ErrorCode: "VALIDATION_ERROR",
			Message:   "Invalid writing ID",
		})
		return
	}

	analysis, err := h.analysisService.GetAnalysis(writingID, userID.(uuid.UUID))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toAnalysisResponse(analysis))
}

func (h *AnalysisHandler) Callback(c *gin.Context) {
	secret := c.GetHeader(CallbackSecretHeader)
	if secret != h.config.MLCallbackSecret {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			ErrorCode: "FORBIDDEN",
			Message:   "Invalid callback secret",
		})
		return
	}

	var req dto.AnalysisCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			ErrorCode: "VALIDATION_ERROR",
			Message:   err.Error(),
		})
		return
	}

	taskID, err := uuid.Parse(req.TaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			ErrorCode: "VALIDATION_ERROR",
			Message:   "Invalid task ID",
		})
		return
	}

	var result *service.CallbackResult
	var callbackErr *service.CallbackError

	if req.Result != nil {
		result = &service.CallbackResult{
			AIProbability: req.Result.AIProbability,
			Feedback:      req.Result.Feedback,
			LatencyMs:     req.Result.LatencyMs,
		}
	}

	if req.Error != nil {
		callbackErr = &service.CallbackError{
			Code:      req.Error.Code,
			Message:   req.Error.Message,
			Retryable: req.Error.Retryable,
		}
	}

	if err := h.analysisService.HandleCallback(taskID, req.Status, result, callbackErr); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "Callback processed"})
}

func toAnalysisResponse(a *model.Analysis) dto.AnalysisResponse {
	resp := dto.AnalysisResponse{
		ID:        a.ID,
		WritingID: a.WritingID,
		Status:    string(a.Status),
		AIScore:   a.AIScore,
		Feedback:  a.Feedback,
		LatencyMs: a.LatencyMs,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}

	if a.ErrorCode != nil {
		code := string(*a.ErrorCode)
		resp.ErrorCode = &code
	}
	resp.ErrorMessage = a.ErrorMessage

	return resp
}
