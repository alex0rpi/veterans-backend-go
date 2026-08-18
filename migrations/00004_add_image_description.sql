-- +goose Up
ALTER TABLE image_metadata
ADD COLUMN file_description TEXT;
-- +goose Down
ALTER TABLE image_metadata DROP COLUMN file_description;