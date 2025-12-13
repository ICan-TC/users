package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Teachers struct {
	bun.BaseModel `bun:"table:teachers,alias:tch"`
	TeacherID     string    `bun:"id,pk"`
	CreatedAt     time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,default:current_timestamp"`
	DeletedAt     time.Time `bun:"deleted_at,default:null"`

	UserID *string `bun:"user_id"`
	User   *Users  `bun:"rel:belongs-to,join:user_id=id"`
}
