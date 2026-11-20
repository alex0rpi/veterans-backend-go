package models

type PositionItem struct {
	ID       int `json:"id" validate:"required"`
	Position int `json:"position" validate:"required,gte=1,lte=100"`
}

type UpdateDisplayPositionsRequest struct {
	NewPositions []PositionItem `json:"new_display_position_items" validate:"required,dive,required,min=1,max=100"`
}
