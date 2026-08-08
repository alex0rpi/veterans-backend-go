-- +goose Up
ALTER TABLE image_metadata
    RENAME COLUMN storage_key TO object_key;
-- +goose Down
ALTER TABLE image_metadata
    RENAME COLUMN object_key TO storage_key;