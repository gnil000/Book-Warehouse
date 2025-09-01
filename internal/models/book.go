package models

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID            uuid.UUID `json:"bookId" binding:"required"`
	DateOfWriting time.Time `json:"year" binding:"required"`
	Title         string    `json:"title" binding:"required,min=1,max=500"`
	Author        Author    `json:"author" binding:"required"`
	Quantity      int       `json:"quantity" binding:"required"`
}

type CreateOrUpdateBookRequest struct {
	DateOfWriting time.Time `json:"year" binding:"required"`
	Title         string    `json:"title" binding:"required,min=1,max=500"`
	Author        Author    `json:"author" binding:"required"`
}

type CreateBookResponse struct {
	ID uuid.UUID `json:"bookId" binding:"required"`
}

type ChangeBookQuantityRequest struct {
	ID       uuid.UUID `json:"bookId" binding:"required"`
	Quantity int       `json:"quantity" binding:"required"`
}

type ChangeBookQuantityResponse struct {
	Quantity int `json:"quantity" binding:"required"`
}

type FindByIdRequest struct {
	ID uuid.UUID `json:"bookId" binding:"required"`
}
