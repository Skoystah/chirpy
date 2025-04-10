package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Token      string
	UserID     uuid.UUID
	Expires_at time.Time
	Revoked_at sql.NullTime
}

type RefreshResponse struct {
	Token string `json:"token"`
}
