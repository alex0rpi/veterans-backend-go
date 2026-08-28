-- +goose Up
ALTER TABLE image_metadata
ADD COLUMN visible BOOLEAN DEFAULT TRUE;
-- +goose Down
ALTER TABLE image_metadata DROP COLUMN visible;