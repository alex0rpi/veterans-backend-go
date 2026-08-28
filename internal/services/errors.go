package services

import "errors"

var (
	ErrUnsupportedMediaType = errors.New(
		"unsupported media type",
	)

	//* File size / dimensions errors

	ErrFileTooLarge = errors.New(
		"file size exceeds the maximum allowed limit",
	)
	ErrImageDimensionsOutOfRange = errors.New(
		"image dimensions are out of the allowed range",
	)
	ErrInvalidMediaContext = errors.New(
		"invalid media context",
	)

	//* Category errors ------------------

	ErrInvalidCategory = errors.New(
		"invalid category",
	)
	ErrCategoryNotAllowedForGivenContext = errors.New(
		"category provided is not allowed for the given media context",
	)
	ErrCategoryRequiredForGivenContext = errors.New(
		"category is required for the given media context",
	)

	//* Season errors ------------------

	ErrInvalidSeason = errors.New(
		"invalid season",
	)
	ErrSeasonRequiredForGivenContext = errors.New(
		"season is required for the given media context",
	)
	ErrSeasonNotAllowedForGivenContext = errors.New(
		"season field is not allowed for the given media context",
	)

	//* Title & description errors
	ErrDocumentTitleRequired = errors.New(
		"document title is required",
	)

	//* Display position errors ------------------

	ErrDisplayPositionRequired = errors.New(
		"display position is required",
	)

	ErrDisplayPositionNotValid = errors.New(
		"display position is not valid",
	)

	ErrDisplayPositionNotAllowedOrUsed = errors.New(
		"display position is either not allowed or already used",
	)

	//* Not found errors ------------------

	ErrMediaNotFound = errors.New(
		"media not found",
	)
)
