-- +goose Up
-- Make display_order column UNIQUE
ALTER TABLE image_metadata
    RENAME COLUMN display_order TO display_position;
ALTER TABLE image_metadata
ALTER COLUMN display_position
SET NOT NULL;
CREATE UNIQUE INDEX unique_display_order_per_group ON image_metadata (
    media_context,
    COALESCE(season, ''),
    COALESCE(category, ''),
    display_position
);
-- +goose Down
DROP INDEX IF EXISTS unique_display_order_per_group;
ALTER TABLE image_metadata
ALTER COLUMN display_position DROP NOT NULL;
ALTER TABLE image_metadata
    RENAME COLUMN display_position TO display_order;