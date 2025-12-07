package models

import (
	"time"

	"github.com/uptrace/bun"
)

type RefreshTokens struct {
	bun.BaseModel `bun:"table:refresh_tokens,alias:rt"`
	ID            string    `bun:"id,pk"`
	UserID        string    `bun:"user_id"`
	Token         string    `bun:"token"`
	Device        string    `bun:"device"`
	ExpiresAt     time.Time `bun:"expires_at"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
}
