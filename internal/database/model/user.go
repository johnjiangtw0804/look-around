package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;uniqueIndex;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserName  string
	Password  string
	Gender    string
}
