package models

import "mime/multipart"

type UploadMediaRequest struct {
	File             multipart.File
	OriginalFilename string
	FileSize         int64
}
