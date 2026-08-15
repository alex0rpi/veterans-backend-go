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
			$11, $12, $13, $14)
		RETURNING id
		`,
		image.ObjectKey,
		image.OriginalFilename,
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
