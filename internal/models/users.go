package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Users struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	ID            string    `bun:"id,pk"`
	Username      string    `bun:"username"`
	Email         string    `bun:"email"`
	PasswordHash  string    `bun:"password_hash"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
}
