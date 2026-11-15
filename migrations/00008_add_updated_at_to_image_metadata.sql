-- +goose Up
ALTER TABLE image_metadata
ADD COLUMN updated_at TIMESTAMPTZ;
-- +goose Down
ALTER TABLE image_metadata DROP COLUMN updated_at;