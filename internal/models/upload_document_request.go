package models

import "mime/multipart"

type UploadDocumentRequest struct {
	File             multipart.File
	OriginalFilename string
	FileSize         int64
	Title            string  `json:"document_title"`
	FileDescription  *string `json:"document_description,omitempty"`
}

// * means that the field is a pointer, which allows it to be nil if not provided.
// * The `omitempty` tag means that if the field is nil, it will be omitted from the JSON output.
