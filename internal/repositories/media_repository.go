package repositories

import (
	"context"
	"veterans-go-chi-server/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MediaRepository struct {
	db *pgxpool.Pool
}

func NewMediaRepository(db *pgxpool.Pool) *MediaRepository {
	return &MediaRepository{
		db: db,
	}
}

func (r *MediaRepository) Create(
	ctx context.Context,
	image *models.ProcessedMedia,
) error {
	// _, err := r.db.Exec(
	return r.db.QueryRow(
		ctx,
		`
		INSERT INTO image_metadata (
			object_key,
			original_filename,
			file_description,
			mime_type,
			width,
			height,
			filesize,
			blur_key,
			small_key,
			medium_key,
			large_key,
			media_context,
			season,
			category,
			display_order
		)
		VALUES ($1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15)
		RETURNING id
		`,
		image.ObjectKey,
		image.OriginalFilename,
		image.FileDescription,
		image.MimeType,
		image.Width,
		image.Height,
		image.FileSize,
		image.BlurKey,
		image.SmallKey,
		image.MediumKey,
		image.LargeKey,
		image.MediaContext,
		image.Season,
		image.Category,
		image.DisplayOrder,
	).Scan(&image.ID)
}

func (r *MediaRepository) Delete(
	ctx context.Context,
	id string,
) error {
	_, err := r.db.Exec(
		ctx,
		"DELETE FROM image_metadata WHERE id = $1",
		id,
	)
	return err
}

func (r *MediaRepository) ListMedia(
	ctx context.Context,
	mediaContext string,
	season *string,
) ([]models.ProcessedMedia, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT id, object_key, original_filename, file_description,
			mime_type, width, height, filesize, blur_key,
			small_key, medium_key, large_key, media_context,
			season, category, display_order
		FROM image_metadata
		WHERE media_context = $1 AND ($2::text IS NULL OR season = $2)
		ORDER BY display_order ASC NULLS LAST, id ASC
		`,
		mediaContext,
		season,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mediaList []models.ProcessedMedia
	for rows.Next() {
		var media models.ProcessedMedia
		if err := rows.Scan(
			&media.ID,
			&media.ObjectKey,
			&media.OriginalFilename,
			&media.FileDescription,
			&media.MimeType,
			&media.Width,
			&media.Height,
			&media.FileSize,
			&media.BlurKey,
			&media.SmallKey,
			&media.MediumKey,
			&media.LargeKey,
			&media.MediaContext,
			&media.Season,
			&media.Category,
			&media.DisplayOrder,
		); err != nil {
			return nil, err
		}
		mediaList = append(mediaList, media)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mediaList, nil
}
