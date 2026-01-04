package dto

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateWritingRequest struct {
	Type    string `json:"type" binding:"required,oneof=essay cover_letter"`
	Title   string `json:"title" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required,max=2000"`
}

type UpdateWritingRequest struct {
	Type    string `json:"type" binding:"omitempty,oneof=essay cover_letter"`
	Title   string `json:"title" binding:"omitempty,min=1,max=255"`
	Content string `json:"content" binding:"omitempty,max=2000"`
}

type ListWritingsQuery struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}
