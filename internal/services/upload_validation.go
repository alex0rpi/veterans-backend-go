package services

// Shared validation helpers reused by mediaService and documentService.

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
