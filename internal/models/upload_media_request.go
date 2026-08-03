package models

import "mime/multipart"

type UploadMediaRequest struct {
	File				multipart.File
	OriginalFilename	string
	MimeType			string
	FileSize			int64
}