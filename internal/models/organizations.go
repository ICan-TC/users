package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Organization struct {
	bun.BaseModel `bun:"table:organizations,alias:org"`

	ID       string                 `bun:"id,pk"` //ulid
	Name     string                 `bun:"name,notnull"`
	Logo     string                 `bun:"logo,notnull"`
	Tags     []string               `bun:"tags,array"`
	Labels   map[string]string      `bun:"labels,type:jsonb,default:'{}'"`
	Metadata map[string]interface{} `bun:"metadata,type:jsonb,default:'{}'"`

	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
	CreatedBy string    `bun:"created_by,notnull"`
}
