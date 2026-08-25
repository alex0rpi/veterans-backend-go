package repositories

import (
	"context"
	"veterans-go-chi-server/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DocumentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{
		db: db,
	}
}

func (r *DocumentRepository) Create(
	ctx context.Context,
	document *models.DocumentMetadata,
) error {
	return r.db.QueryRow(
		ctx,
		`
		INSERT INTO document_metadata (
			object_key,
			original_filename,
			mime_type,
			title,
			description,
			filesize
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
		`,
		document.ObjectKey,
		document.OriginalFilename,
		document.MimeType,
		document.Title,
		document.Description,
		document.FileSize,
	).Scan(&document.ID, &document.CreatedAt, &document.UpdatedAt)
}

func (r *DocumentRepository) List(
	ctx context.Context,
) ([]*models.DocumentMetadata, error) {
	rows, err := r.db.Query(
		ctx,
		"SELECT id, object_key, original_filename, mime_type, title, description, filesize, created_at, updated_at FROM document_metadata ORDER BY created_at",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := make([]*models.DocumentMetadata, 0)
	for rows.Next() {
		document := &models.DocumentMetadata{}
		if err := rows.Scan(
			&document.ID,
			&document.ObjectKey,
			&document.OriginalFilename,
			&document.MimeType,
			&document.Title,
			&document.Description,
			&document.FileSize,
			&document.CreatedAt,
			&document.UpdatedAt,
		); err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *DocumentRepository) Delete(
	ctx context.Context,
	id string,
) error {
	_, err := r.db.Exec(
		ctx,
		"DELETE FROM document_metadata WHERE id = $1",
		id,
	)
	return err
}
