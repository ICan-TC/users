package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Enrollments struct {
	bun.BaseModel `bun:"table:enrollments,alias:enr"`
	StudentID     string    `bun:"student_id,pk"`
	GroupID       string    `bun:"group_id,pk"`
	Fee           float64   `bun:"fee"`
	CreatedAt     time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,default:current_timestamp"`
	DeletedAt     time.Time `bun:"deleted_at,default:null"`

	Student *Students `bun:"rel:belongs-to,join:student_id=id"`
	Group   *Groups   `bun:"rel:belongs-to,join:group_id=id"`
}
