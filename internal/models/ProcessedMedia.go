package models

type ProcessedMedia struct {
	ID int64
	ObjectKey string
	OriginalFilename string
	MimeType         string
	FileSize         int64
	Width  int
	Height int
}
