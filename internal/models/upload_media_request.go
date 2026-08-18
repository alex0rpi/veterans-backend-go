package models

import "mime/multipart"

type UploadMediaRequest struct {
	File             multipart.File
	OriginalFilename string
	FileSize         int64

	FileDescription *string `json:"file_description,omitempty"`
	MediaContext    *string `json:"media_context,omitempty"`
	Season          *string `json:"season,omitempty"`
	Category        *string `json:"category,omitempty"`
	DisplayOrder    *int    `json:"display_order,omitempty"`
}
