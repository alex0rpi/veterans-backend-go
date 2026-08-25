-- +goose Up
CREATE TABLE document_metadata (
    id BIGSERIAL PRIMARY KEY,
    object_key TEXT NOT NULL UNIQUE,
    original_filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    filesize BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose Down
DROP TABLE document_metadata;