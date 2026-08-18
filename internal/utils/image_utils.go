package utils

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
)

func GetImageDimensions(file multipart.File) (int, int, error) {

	cfg, _, err := image.DecodeConfig(file)

	if err != nil {
		return 0, 0, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return 0, 0, err
	}

	return cfg.Width, cfg.Height, nil
}

func GetMimeType(file multipart.File) (string, error) {

	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

func ImageDimensionsAreValid(width, height int) bool {
	const maxDimension = 8000

	return width > 0 &&
		width <= maxDimension &&
		height > 0 &&
		height <= maxDimension
}
