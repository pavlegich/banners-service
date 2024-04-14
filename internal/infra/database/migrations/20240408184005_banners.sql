-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS banners (
    id serial PRIMARY KEY,
    tag_ids integer[] NOT NULL,
    feature_id integer NOT NULL,
    content jsonb,
    is_active boolean,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

-- create indexes
CREATE INDEX IF NOT EXISTS tag_ids_idx ON banners (tag_ids);
CREATE INDEX IF NOT EXISTS feature_id_idx ON banners (feature_id);
CREATE INDEX IF NOT EXISTS updated_at_idx ON banners (updated_at);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP INDEX updated_at_idx;
DROP INDEX feature_id_idx;
DROP INDEX tag_ids_idx;
DROP TABLE banners;
