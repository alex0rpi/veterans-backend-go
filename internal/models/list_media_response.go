package models

type ListMediaResponse struct {
	OriginalFilename string		`json:"original_filename"`
	FileDescription  *string	`json:"file_description"`
	Width            int		`json:"width"`
	Height           int		`json:"height"`

	BlurMediaURL   string		`json:"blur_url"`
	SmallMediaURL  string		`json:"small_url"`
	MediumMediaURL string		`json:"medium_url"`
	LargeMediaURL  string		`json:"large_url"`

	MediaContext *string	`json:"media_context,omitempty"`
	Season       *string	`json:"season,omitempty"`
	Category     *string	`json:"category,omitempty"`
	DisplayOrder *int		`json:"display_order,omitempty"`
}
