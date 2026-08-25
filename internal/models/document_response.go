package models

type DocumentResponse struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	URL         string  `json:"url"`
	MimeType    string  `json:"mime_type"`
	FileSize    int64   `json:"filesize"`
}
