package models

import "github.com/google/uuid"

type Author struct {
	ID         uuid.UUID
	FirstName  string
	SecondName string
	Surname    string
}
