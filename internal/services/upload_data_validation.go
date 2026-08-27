package services

// Validation helpers specific to the mediaService.

import (
	"strconv"
	"time"
	"veterans-go-chi-server/internal/constants"
	"veterans-go-chi-server/internal/models"
)

func validateUploadMetadata(request models.UploadMediaRequest) error {
	if !isValidMediaContext(request.MediaContext) {
		return ErrInvalidMediaContext
	}

	if request.Category != nil && !isValidCategory(*request.Category) {
		return ErrInvalidCategory
	}

	if request.Season != nil && !isValidSeason(*request.Season) {
		return ErrInvalidSeason
	}

	return nil
}

func validateListMediaRequest(request models.ListMediaRequest) error {

	if !isValidMediaContext(request.MediaContext) {
		return ErrInvalidMediaContext
	}

	if request.Season != nil && !isValidSeason(*request.Season) {
		return ErrInvalidSeason
	}

	if request.MediaContext == string(constants.MediaContextSeasonSlider) && request.Season == nil {
		return ErrSeasonRequiredForGivenContext
	}

	if request.MediaContext != string(constants.MediaContextSeasonSlider) && request.Season != nil {
		return ErrSeasonNotAllowedForGivenContext
	}

	return nil
}

func isValidMediaContext(value string) bool {
	switch value {
	case
		string(constants.MediaContextSeasonSlider),
		string(constants.MediaContextHomeSlider),
		string(constants.MediaContextSingle):
		return true
	default:
		return false
	}
}

func isValidCategory(value string) bool {
	switch value {
	case
		string(constants.ProCategory),
		string(constants.BasesCategory),
		string(constants.VetsCategory),
		string(constants.BoardCategory),
		string(constants.HistoricBoardCategory):
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

	currentYear := time.Now().Year() % 100 // Get the last two digits of the current year
	if start < 0 || start > currentYear || end < 0 {
		return false
	}

	return end == (start+1)%100
}
