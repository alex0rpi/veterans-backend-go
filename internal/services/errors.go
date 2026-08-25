package services

import "errors"

var (
	ErrUnsupportedMediaType = errors.New(
		"unsupported media type",
	)

	ErrFileTooLarge = errors.New(
		"file size exceeds the maximum allowed limit",
	)

	ErrImageDimensionsOutOfRange = errors.New(
		"image dimensions are out of the allowed range",
	)

	ErrInvalidMediaContext = errors.New(
		"invalid media context",
	)

	ErrInvalidMediaCategory = errors.New(
		"invalid media category",
	)

	ErrInvalidSeason = errors.New(
		"invalid season",
	)

	ErrMediaNotFound = errors.New(
		"media not found",
	)

	ErrDocumentTitleRequired = errors.New(
		"document title is required",
	)

	ErrSeasonRequiredForSeasonSlider = errors.New(
		"season is required when media_context is 'season_slider'",
	)

	ErrSeasonNotAllowedForRequestedContext = errors.New(
		"season field is not allowed for the requested media context",
	)
)
