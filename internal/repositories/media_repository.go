package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	mediamodel "veterans-go-chi-server/internal/models/media"

	"github.com/jackc/pgx/v5"
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
	image *mediamodel.ProcessedMedia,
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
			display_position,
			visible
			
		)
		VALUES ($1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15,
			$16)
		RETURNING id, created_at, updated_at
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
		image.DisplayPosition,
		image.Visible,
	).Scan(&image.ID, &image.CreatedAt, &image.UpdatedAt)
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

func (r *MediaRepository) GetImageMetadataById(
	ctx context.Context,
	id int,
) (*mediamodel.ImageMetadata, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT id, object_key, original_filename, file_description, mime_type, width, height, filesize, blur_key, small_key, medium_key, large_key, media_context, season, category, display_position, visible
		FROM image_metadata
		WHERE id = $1
		`,
		id,
	)

	var image mediamodel.ImageMetadata
	err := row.Scan(
		&image.ID,
		&image.ObjectKey,
		&image.ImageInfo.OriginalFilename,
		&image.ImageInfo.FileDescription,
		&image.ImageMetadata.MimeType,
		&image.ImageMetadata.Width,
		&image.ImageMetadata.Height,
		&image.ImageMetadata.FileSize,
		&image.URLs.Blur,
		&image.URLs.Small,
		&image.URLs.Medium,
		&image.URLs.Large,
		&image.ImageInfo.MediaContext,
		&image.ImageInfo.Season,
		&image.ImageInfo.Category,
		&image.ImageInfo.DisplayPosition,
		&image.ImageInfo.Visible,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("image with given id not found")
		}
		return nil, err
	}
	return &image, nil
}

func (r *MediaRepository) ListMedia(
	ctx context.Context,
	mediaContext string,
	season *string,
	category *string,
) ([]mediamodel.ProcessedMedia, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT object_key, original_filename, file_description,
			mime_type, width, height, filesize, blur_key,
			small_key, medium_key, large_key, media_context,
			season, category, display_position, visible, created_at, updated_at
		FROM image_metadata
		WHERE media_context = $1 AND ($2::text IS NULL OR season = $2) AND ($3::text IS NULL OR category = $3)
		ORDER BY display_position ASC NULLS LAST, id ASC
		`,
		mediaContext,
		season,
		category,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mediaList []mediamodel.ProcessedMedia
	for rows.Next() {
		var media mediamodel.ProcessedMedia
		if err := rows.Scan(
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
			&media.DisplayPosition,
			&media.Visible,
			&media.CreatedAt,
			&media.UpdatedAt,
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

func (r *MediaRepository) UpdateDisplayPositions(
	ctx context.Context,
	newPositions []mediamodel.PositionItem,
) error {
	// Usar batch para actualizar múltiples registros de forma eficiente
	batch := &pgx.Batch{}

	for _, item := range newPositions {
		batch.Queue(
			`UPDATE image_metadata SET display_position = $1 WHERE id = $2`,
			item.Position,
			item.ID,
		)
	}

	// Enviar todas las queries en una transacción
	results := r.db.SendBatch(ctx, batch)
	defer func() { _ = results.Close() }()

	// Comprobar que todas las updates se ejecutaron
	for i := 0; i < batch.Len(); i++ {
		commandTag, err := results.Exec()
		if err != nil {
			return err
		}
		if commandTag.RowsAffected() == 0 {
			return fmt.Errorf("image id %d not found", newPositions[i].ID)
		}
		log.Printf("Updated %d rows for position item %d", commandTag.RowsAffected(), i)
	}

	return nil
}

func (r *MediaRepository) UpdateMedia(
	ctx context.Context,
	id int,
	updateData mediamodel.UpdateMediaItemRequest,
) (*mediamodel.UpdatedMediaItem, error) {
	row := r.db.QueryRow(
		ctx,
		`UPDATE image_metadata
		SET original_filename = COALESCE($2, original_filename),
			file_description = COALESCE($3, file_description),
			media_context = COALESCE($4, media_context),
			season = COALESCE($5, season),
			category = COALESCE($6, category),
			visible = COALESCE($7, visible),
			updated_at = NOW()
		WHERE id = $1
		RETURNING original_filename, file_description, media_context, season, category, visible, updated_at`,
		id,
		updateData.OriginalFilename,
		updateData.FileDescription,
		updateData.MediaContext,
		updateData.Season,
		updateData.Category,
		updateData.Visible,
	)
	var updated mediamodel.UpdatedMediaItem
	err := row.Scan(
		&updated.OriginalFilename,
		&updated.FileDescription,
		&updated.MediaContext,
		&updated.Season,
		&updated.Category,
		&updated.Visible,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *MediaRepository) GetImageByID(ctx context.Context, id int) (*mediamodel.ProcessedMedia, error) {
	row := r.db.QueryRow(
		ctx,
		`
		SELECT id, object_key
		FROM image_metadata
		WHERE id = $1
		`,
		id,
	)

	var media mediamodel.ProcessedMedia
	if err := row.Scan(
		&media.ID,
		&media.ObjectKey,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return &media, nil
}

func (r *MediaRepository) GetUsedDisplayPositions(ctx context.Context, mediaContext string, season, category *string) ([]int, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT display_position
		FROM image_metadata
		WHERE media_context = $1
		AND COALESCE(season, '') = COALESCE($2, '')
		AND COALESCE(category, '') = COALESCE($3, '')
		ORDER BY display_position ASC
		`,
		mediaContext,
		season,
		category,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usedDisplayPositions []int
	for rows.Next() {
		var displayOrder int
		if err := rows.Scan(&displayOrder); err != nil {
			return nil, err
		}
		usedDisplayPositions = append(usedDisplayPositions, displayOrder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return usedDisplayPositions, nil

}

func (r *MediaRepository) DeleteImage(ctx context.Context, id int) error {
	_, err := r.db.Exec(
		ctx,
		"DELETE FROM image_metadata WHERE id = $1",
		id,
	)
	return err
}
