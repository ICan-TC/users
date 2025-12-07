CREATE TABLE IF NOT EXISTS organizations (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	logo TEXT NOT NULL,
	tags TEXT[],
	labels JSONB DEFAULT '{}'::JSONB,
	metadata JSONB DEFAULT '{}'::JSONB,

	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	created_by TEXT NOT NULL DEFAULT 'system'
);

CREATE INDEX IF NOT EXISTS idx_organizations_name ON organizations(name);
CREATE INDEX IF NOT EXISTS idx_organizations_tags ON organizations USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_organizations_labels ON organizations USING GIN(labels);
