package media

type ListMediaRequest struct {
	MediaContext string  `json:"media_context,omitempty"`
	Season       *string `json:"season,omitempty"`
	Category     *string `json:"category,omitempty"`
}
