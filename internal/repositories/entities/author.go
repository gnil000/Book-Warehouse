package entities

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key"`
	DateOfBirth time.Time `gorm:"type:timestampz"`
}
