package schema

import "github.com/google/uuid"

type UserLike struct {
	UserUD  uuid.UUID `gorm:"type:uuid;primary_key;uniqueIndex;default:uuid_generate_v4()"`
	EventID uuid.UUID `gorm:"type:uuid;primary_key;uniqueIndex;default:uuid_generate_v4()"`
}
