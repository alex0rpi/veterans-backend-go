package media

import "time"

// Individual image item
type UpdatedMediaItem struct {
	// Metadata
	OriginalFilename string  `json:"original_filename"`
	FileDescription  *string `json:"file_description"`

	MediaContext string     `json:"media_context"`
	Season       *string    `json:"season,omitempty"`
	Category     *string    `json:"category,omitempty"`
	Visible      bool       `json:"visible"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
