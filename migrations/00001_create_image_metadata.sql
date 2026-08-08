-- +goose Up
CREATE TABLE image_metadata (
    id BIGSERIAL PRIMARY KEY,
    storage_key TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    filesize BIGINT NOT NULL,
    blur_key TEXT,
    small_key TEXT,
    medium_key TEXT,
    large_key TEXT
);
-- +goose Down
DROP TABLE IF EXISTS image_metadata;