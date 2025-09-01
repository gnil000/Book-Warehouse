package entities

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	DateOfBirth time.Time `gorm:"type:date"`
	FirstName   string    `gorm:"type:text"`
	SecondName  string    `gorm:"type:text"`
	Surname     string    `gorm:"type:text"`
}
