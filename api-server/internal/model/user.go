package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email            string     `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash     string     `gorm:"not null;size:255" json:"-"`
	DailySubmitCount int        `gorm:"not null;default:0" json:"daily_submit_count"`
	LastSubmitDate   *time.Time `gorm:"type:date" json:"last_submit_date"`
	CreatedAt        time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"not null;default:now()" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
