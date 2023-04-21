package schema

import "github.com/google/uuid"

type UserLikeGenre struct {
	UserID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Genre    string    `gorm:"type:varchar(255)"`
	SubGenre string    `gorm:"type:varchar(255)"`
}
