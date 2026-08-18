package models

type ProcessedMedia struct {
	ID int64

	ObjectKey        string
	OriginalFilename string
	FileDescription  *string
	MimeType         string
	FileSize         int64
	Width            int
	Height           int

	BlurKey   string
	SmallKey  string
	MediumKey string
	LargeKey  string

	MediaContext *string
	Season       *string
	Category     *string
	DisplayOrder *int
}
