package models

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID          uuid.UUID `json:"author_id" binding:"required"`
	DateOfBirth time.Time `json:"birth" binding:"required"`
	FirstName   string    `json:"first_name" binding:"required"`
	SecondName  string    `json:"second_name" binding:"required"`
	Surname     string    `json:"surname" binding:"required"`
}

type CreateAuthorRequest struct {
	DateOfBirth time.Time `json:"birth" binding:"required"`
	FirstName   string    `json:"first_name" binding:"required"`
	SecondName  string    `json:"second_name" binding:"required"`
	Surname     string    `json:"surname" binding:"required"`
}

type CreateAuthorResponse struct {
	ID uuid.UUID `json:"author_id" binding:"required"`
}

type UpdateAuthorRequest struct {
	ID          uuid.UUID `json:"author_id" binding:"required"`
	DateOfBirth time.Time `json:"birth"`
	FirstName   string    `json:"first_name"`
	SecondName  string    `json:"second_name"`
	Surname     string    `json:"surname"`
}
