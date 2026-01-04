package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/truegul/api-server/internal/data"
	"github.com/truegul/api-server/internal/dto"
	"github.com/truegul/api-server/internal/service"
)

type WritingHandler struct {
	writingService *service.WritingService
}

func NewWritingHandler(writingService *service.WritingService) *WritingHandler {
	return &WritingHandler{writingService: writingService}
}

func (h *WritingHandler) Create(c *gin.Context) {
	var req dto.CreateWritingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err.Error())
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	writing, err := h.writingService.Create(userID, req.Type, req.Title, req.Content)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toWritingResponse(writing))
}

func (h *WritingHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleValidationError(c, "Invalid writing ID")
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	writing, err := h.writingService.GetByID(id, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toWritingResponse(writing))
}

func (h *WritingHandler) List(c *gin.Context) {
	var query dto.ListWritingsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		handleValidationError(c, err.Error())
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	writings, total, err := h.writingService.List(userID, query.Page, query.Limit)
	if err != nil {
		handleError(c, err)
		return
	}

	writingResponses := make([]dto.WritingResponse, len(writings))
	for i, w := range writings {
		writingResponses[i] = toWritingResponse(w)
	}

	totalPages := int(total) / query.Limit
	if int(total)%query.Limit != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, dto.WritingListResponse{
		Writings:   writingResponses,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPages,
	})
}

func (h *WritingHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleValidationError(c, "Invalid writing ID")
		return
	}

	var req dto.UpdateWritingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleValidationError(c, err.Error())
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	var writingType, title, content *string
	if req.Type != "" {
		writingType = &req.Type
	}
	if req.Title != "" {
		title = &req.Title
	}
	if req.Content != "" {
		content = &req.Content
	}

	writing, err := h.writingService.Update(id, userID, writingType, title, content)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toWritingResponse(writing))
}

func (h *WritingHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		handleValidationError(c, "Invalid writing ID")
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	err = h.writingService.Delete(id, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "Writing deleted successfully",
	})
}

func toWritingResponse(w *data.Writing) dto.WritingResponse {
	return dto.WritingResponse{
		ID:          w.ID,
		UserID:      w.UserID,
		Type:        string(w.Type),
		Title:       w.Title,
		Content:     w.Content,
		Status:      string(w.Status),
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
		SubmittedAt: w.SubmittedAt,
	}
}
