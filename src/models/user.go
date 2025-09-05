package models

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID `json:"-"`
	Login          string    `json:"login"`
	Password       string    `json:"password"`
	HashedPassword string    `json:"-"`
}

type RegistrationUserRequest struct {
	Login    string `json:"login" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=8"`
}
