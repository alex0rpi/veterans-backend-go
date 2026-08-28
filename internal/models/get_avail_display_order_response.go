package models

type GetAvailDisplayPositionsResponse struct {
	NextAvailableDisplayPosition int   `json:"next_avail_display_position"`
	UnusedDisplayPositions       []int `json:"unused_middle_display_positions"`
}
