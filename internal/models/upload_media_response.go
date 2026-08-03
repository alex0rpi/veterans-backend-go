package models

type UploadMediaResponse struct {
	Filename string `json:"filename"`
	Size	 int64  `json:"size"`
	MimeType string `json:"mimeType"`
}