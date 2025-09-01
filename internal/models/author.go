package models

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID          uuid.UUID
	DateOfBirth time.Time
	FirstName   string
	SecondName  string
	Surname     string
	FullName    string
}
