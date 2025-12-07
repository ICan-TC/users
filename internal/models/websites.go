package models

import (
	"time"

	"github.com/uptrace/bun"
)

// Website represents a monitored website entity.
// Matches the SQL schema in migrations/01_create_websites_table.up.sql.
type Website struct {
	bun.BaseModel `bun:"table:websites,alias:w"`

	ID          string                 `bun:"id,pk"`
	Name        string                 `bun:"name,notnull"`
	URL         string                 `bun:"url,notnull"`
	Environment string                 `bun:"environment,notnull,default:'production'"`
	Description *string                `bun:"description"`
	Tags        []string               `bun:"tags,array"`
	Labels      map[string]string      `bun:"labels,type:jsonb,default:'{}'"`
	Metadata    map[string]interface{} `bun:"metadata,type:jsonb,default:'{}'"`

	ProjectID string `bun:"project_id,notnull"`

	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
	CreatedBy string    `bun:"created_by,notnull,default:'system'"`
}
