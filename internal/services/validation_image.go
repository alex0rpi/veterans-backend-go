package services

import (
	"veterans-go-chi-server/internal/models"
	"veterans-go-chi-server/internal/utils"
)

// Shared validation helpers reused by mediaService and documentService.

func validateImageProperties(request models.UploadMediaRequest) (string, int, int, error) {
	if err := validateFileSize(request.FileSize, maxImageFileSize); err != nil {
		return "", 0, 0, err
	}

	mimeType, err := utils.GetMimeType(request.File)
	if err != nil {
		return "", 0, 0, err
	}

	if err := validateMimeType(mimeType, allowedImageMimeTypes); err != nil {
		return "", 0, 0, err
	}

	width, height, err := utils.GetImageDimensions(request.File)
	if err != nil {
		return "", 0, 0, err
	}
	if !utils.ImageDimensionsAreValid(width, height) {
		return "", 0, 0, ErrImageDimensionsOutOfRange
	}

	return mimeType, width, height, nil
}

func validateFileSize(fileSize, maxFileSize int64) error {
	if fileSize > maxFileSize {
		return ErrFileTooLarge
	}
	return nil
}

func validateMimeType(mimeType string, allowed map[string]bool) error {
	if !allowed[mimeType] {
		return ErrUnsupportedMediaType
	}
	return nil
}
