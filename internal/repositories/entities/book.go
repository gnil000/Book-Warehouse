package entities

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	DateOfWriting time.Time `gorm:"type:date"`
	Title         string    `gorm:"type:text"`
	AuthorID      uuid.UUID `gorm:"type:uuid"`
	Quantity      int       `gorm:"type:int"`
	Author        Author
}
