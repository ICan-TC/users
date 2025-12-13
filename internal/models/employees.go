package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Employees struct {
	bun.BaseModel `bun:"table:employees,alias:emp"`
	EmployeeID    string    `bun:"id,pk"`
	Role          string    `bun:"role"`
	Salary        float64   `bun:"salary"`
	CreatedAt     time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,default:current_timestamp"`
	DeletedAt     time.Time `bun:"deleted_at,default:null"`

	UserID string `bun:"user_id"`
	User   *Users `bun:"rel:belongs-to,join:user_id=id"`
}
