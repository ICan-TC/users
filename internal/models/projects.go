package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Project struct {
	bun.BaseModel `bun:"table:projects,alias:pr"`

	ID          string                 `bun:"id,pk"`
	Name        string                 `bun:"name,notnull"`
	Description *string                `bun:"description"`
	Tags        []string               `bun:"tags,array"`
	Labels      map[string]string      `bun:"labels,type:jsonb,default:'{}'"`
	Metadata    map[string]interface{} `bun:"metadata,type:jsonb,default:'{}'"`

	OrganizationID string `bun:"organization_id,notnull"`

	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
	CreatedBy string    `bun:"created_by,notnull"`
}
