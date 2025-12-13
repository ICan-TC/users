package models

import (
	"time"

	"github.com/uptrace/bun"
)

type StudentParents struct {
	bun.BaseModel `bun:"table:student_parents,alias:sp"`
	StudentID     string    `bun:"student_id,pk"`
	ParentID      string    `bun:"parent_id,pk"`
	CreatedAt     time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,default:current_timestamp"`
	DeletedAt     time.Time `bun:"deleted_at,default:null"`

	Student *Students `bun:"rel:belongs-to,join:student_id=id"`
	Parent  *Parents  `bun:"rel:belongs-to,join:parent_id=id"`
}

type Parents struct {
	bun.BaseModel `bun:"table:parents,alias:par"`
	ParentID      string    `bun:"id,pk"`
	UserID        string    `bun:"user_id"`
	CreatedAt     time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,default:current_timestamp"`
	DeletedAt     time.Time `bun:"deleted_at,default:null"`

	User     *Users      `bun:"rel:belongs-to,join:user_id=id"`
	Students []*Students `bun:"m2m:student_parents,join:Parent=Student"`
}
