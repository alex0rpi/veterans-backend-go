package models

import "time"

type DocumentMetadata struct {
	ID               int64
	ObjectKey        string
	OriginalFilename string
	MimeType         string
	FileSize         int64
	Title            string
	Description      *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
