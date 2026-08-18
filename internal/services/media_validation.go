package services

import (
	"strconv"
	"veterans-go-chi-server/internal/constants"
	"veterans-go-chi-server/internal/models"
)

func validateUploadMetadata(request models.UploadMediaRequest) error {
	if request.MediaContext != nil && !isValidMediaContext(*request.MediaContext) {
		return ErrInvalidMediaContext
	}

	if request.Category != nil && !isValidMediaCategory(*request.Category) {
		return ErrInvalidMediaCategory
	}

	if request.Season != nil && !isValidSeason(*request.Season) {
		return ErrInvalidSeason
	}

	return nil
}

func isValidMediaContext(value string) bool {
	switch value {
	case string(constants.MediaContextSeasonSlider), string(constants.MediaContextHomeSlider):
		return true
	default:
		return false
	}
}

func isValidMediaCategory(value string) bool {
	switch value {
	case string(constants.ProCategory), string(constants.BasesCategory), string(constants.VetsCategory), string(constants.BoardCategory), string(constants.OtherCategory):
		return true
	default:
		return false
	}
}

func isValidSeason(value string) bool {
	if len(value) != 4 {
		return false
	}
	// strconv is a package in Go that provides functions for converting strings to other types, such as integers.
	// Atoi is a method that converts a string to an integer. It returns the integer value and an error if the conversion fails.
	// :2 means take the first two characters of the string, and [2:] means take the last two characters of the string.
	start, err := strconv.Atoi(value[:2])
	if err != nil {
		return false
	}

	end, err := strconv.Atoi(value[2:])
	if err != nil {
		return false
	}

	return end == (start+1)%100
}
