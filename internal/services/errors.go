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
)