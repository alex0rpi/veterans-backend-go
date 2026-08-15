-- +goose Up
ALTER TABLE image_metadata
ADD COLUMN media_context TEXT,
    ADD COLUMN season TEXT,
    ADD COLUMN category TEXT,
    ADD COLUMN display_order INTEGER;
-- +goose Down
ALTER TABLE image_metadata DROP COLUMN media_context,
    DROP COLUMN season,
    DROP COLUMN category,
    DROP COLUMN display_order;