package models

type GetAvailDisplayPositionsRequest struct {
	MediaContext string  `json:"media_context"`
	Season       *string `json:"season"`
	Category     *string `json:"category"`
}
