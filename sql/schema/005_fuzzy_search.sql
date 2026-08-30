-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create GIN trigram indexes for fast similarity searches
CREATE INDEX IF NOT EXISTS posts_title_trgm_idx ON posts USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS posts_description_trgm_idx ON posts USING gin (description gin_trgm_ops);

-- +goose Down
DROP INDEX IF EXISTS posts_description_trgm_idx;
DROP INDEX IF EXISTS posts_title_trgm_idx;
DROP EXTENSION IF EXISTS pg_trgm;