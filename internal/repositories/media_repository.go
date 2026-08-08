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
			filesize
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
		`,
		image.ObjectKey,
		image.OriginalFilename,
		image.MimeType,
		image.Width,
		image.Height,
		image.FileSize,
	).Scan(&image.ID)
}
