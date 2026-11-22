package models

import "time"

type ImageMetadata struct {
	ID        int64  `json:"id"`
	ObjectKey string `json:"object_key"`

	ImageInfo struct {
		FileDescription  *string `json:"file_description"`
		OriginalFilename string  `json:"original_filename"`
		MediaContext     string  `json:"media_context"`
		Season           *string `json:"season"`
		Category         *string `json:"category"`
		DisplayPosition  *int    `json:"display_position"`
		Visible          bool    `json:"visible"`
	} `json:"image_info"`

	ImageMetadata struct {
		MimeType  string    `json:"mime_type"`
		Width     int       `json:"width"`
		Height    int       `json:"height"`
		FileSize  int64     `json:"file_size"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"image_metadata"`

	// URLs for different media sizes
	URLs struct {
		Blur   string `json:"blur_url"`   // miniatura placeholder
		Small  string `json:"small_url"`  // Mobile
		Medium string `json:"medium_url"` // Tablet
		Large  string `json:"large_url"`
	} `json:"urls"`
}
