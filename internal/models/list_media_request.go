package models

type ListMediaRequest struct {
	MediaContext string  `json:"media_context,omitempty"`
	Season       *string `json:"season,omitempty"`
}
