package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Groups struct {
	bun.BaseModel `bun:"table:groups,alias:grp"`
	GroupID       string                 `bun:"id,pk"`
	Name          string                 `bun:"name"`
	Description   string                 `bun:"description"`
	TeacherID     string                 `bun:"teacher_id"`
	DefaultFee    float64                `bun:"default_fee"`
	Subject       string                 `bun:"subject"`
	Level         string                 `bun:"level"`
	Metadata      map[string]interface{} `bun:"metadata,type:jsonb"`
	CreatedAt     time.Time              `bun:"created_at,default:current_timestamp"`
	UpdatedAt     time.Time              `bun:"updated_at,default:current_timestamp"`
	DeletedAt     time.Time              `bun:"deleted_at,default:null"`

	Teacher *Teachers `bun:"rel:belongs-to,join:teacher_id=id"`
}
