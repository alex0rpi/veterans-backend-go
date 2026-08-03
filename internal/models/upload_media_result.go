package models

type UploadMediaResult struct {
	OriginalFilename	string
	MimeType			string
	FileSize			int64
	Width				int
	Height				int
}