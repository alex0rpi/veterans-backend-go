package models

import "time"

type ImageMetadata struct {
	ID               int64
	ObjectKey        string
	OriginalFilename string
	MimeType         string
	CreatedAt        time.Time
	Width            int
	Height           int
	FileSize         int64

	BlurKey   *string
	SmallKey  *string
	MediumKey *string
	LargeKey  *string
}
