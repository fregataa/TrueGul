package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/truegul/api-server/internal/data"
	apperrors "github.com/truegul/api-server/internal/errors"
	"github.com/truegul/api-server/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *data.User) error {
	m := toModel(user)
	if err := r.db.Create(m).Error; err != nil {
		return apperrors.InternalServerWrap(err, "Failed to create user")
	}
	user.ID = m.ID
	user.CreatedAt = m.CreatedAt
	user.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*data.User, error) {
	var m model.User
	err := r.db.Where("email = ?", email).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("User not found")
		}
		return nil, apperrors.InternalServerWrap(err, "Failed to find user")
	}
	return toData(&m), nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*data.User, error) {
	var m model.User
	err := r.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("User not found")
		}
		return nil, apperrors.InternalServerWrap(err, "Failed to find user")
	}
	return toData(&m), nil
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, apperrors.InternalServerWrap(err, "Failed to check user existence")
	}
	return count > 0, nil
}

func (r *UserRepository) Update(user *data.User) error {
	m := toModel(user)
	if err := r.db.Save(m).Error; err != nil {
		return apperrors.InternalServerWrap(err, "Failed to update user")
	}
	user.UpdatedAt = m.UpdatedAt
	return nil
}

func toModel(d *data.User) *model.User {
	return &model.User{
		ID:               d.ID,
		Email:            d.Email,
		PasswordHash:     d.PasswordHash,
		DailySubmitCount: d.DailySubmitCount,
		LastSubmitDate:   d.LastSubmitDate,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}

func toData(m *model.User) *data.User {
	return &data.User{
		ID:               m.ID,
		Email:            m.Email,
		PasswordHash:     m.PasswordHash,
		DailySubmitCount: m.DailySubmitCount,
		LastSubmitDate:   m.LastSubmitDate,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}
