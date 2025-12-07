package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Users struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	UserID        string     `bun:"id,pk"`
	Username      string     `bun:"username"`
	Email         string     `bun:"email"`
	PasswordHash  string     `bun:"password_hash"`
	FirstName     *string    `bun:"first_name"`
	FamilyName    *string    `bun:"family_name"`
	PhoneNumber   *string    `bun:"phone_number"`
	DateOfBirth   *time.Time `bun:"date_of_birth"`
	CreatedAt     time.Time  `bun:"created_at,default:current_timestamp"`
	UpdatedAt     time.Time  `bun:"updated_at,default:current_timestamp"`
}
