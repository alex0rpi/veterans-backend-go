package media

// Individual image item
type MediaItem struct {
	FileDescription *string `json:"file_description"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`

	Season          *string `json:"season,omitempty"`
	Category        *string `json:"category,omitempty"`
	DisplayPosition int     `json:"display_position"`
	Visible         bool    `json:"visible"`

	// URLs for different media sizes
	URLs struct {
		Blur   string `json:"blur_url"`   // miniatura placeholder
		Small  string `json:"small_url"`  // Mobile
		Medium string `json:"medium_url"` // Tablet
		Large  string `json:"large_url"`
	} `json:"urls"`
}

// Response for listing media items
type ListMediaResponse struct {
	MediaItems []MediaItem `json:"media_items"`
	Total      int         `json:"total"`
}
