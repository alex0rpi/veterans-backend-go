package models

type ProcessedMedia struct {
	ID int64

	ObjectKey        string  `json:"object_key"`
	OriginalFilename string  `json:"original_filename"`
	FileDescription  *string `json:"file_description"`
	MimeType         string  `json:"mime_type"`
	FileSize         int64   `json:"file_size"`
	Width            int     `json:"width"`
	Height           int     `json:"height"`

	BlurKey   string `json:"blur_key"`
	SmallKey  string `json:"small_key"`
	MediumKey string `json:"medium_key"`
	LargeKey  string `json:"large_key"`

	MediaContext    string  `json:"media_context"`
	Season          *string `json:"season,omitempty"`
	Category        *string `json:"category,omitempty"`
	DisplayPosition int     `json:"display_position"`
	Visible         *bool   `json:"visible"`
}
