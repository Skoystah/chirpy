package model

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ExpiresIn int    `json:"expires_in_seconds"`
}

type LoginError struct {
	Error string `json:"error"`
}

type LoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}
