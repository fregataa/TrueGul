package data

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID
	Email            string
	PasswordHash     string
	DailySubmitCount int
	LastSubmitDate   *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
