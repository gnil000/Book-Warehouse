package entities

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key"`
	DateOfWriting time.Time `gorm:"type:timestampz"`
	Title         string    `gorm:"type:text"`
	AuthorID      uuid.UUID `gorm:"type:uuid"`
	Quantity      int       `gorm:"type:int"`
	Author        Author
}
