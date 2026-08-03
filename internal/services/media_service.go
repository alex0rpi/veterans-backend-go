package services

import (
	"context"
	"image"
	"io"
	"veterans-go-chi-server/internal/models"
)

// This interface is a contract that defines the methods that any media service implementation must provide.
type MediaService interface {
	Upload(
		ctx context.Context,
		request models.UploadMediaRequest,
	) (*models.UploadMediaResult, error)
}

type mediaService struct {
}

// This function is effectively a constructor. It returns a pointer to a mediaService struct, which implements the MediaService interface.
func NewMediaService() MediaService {
	return &mediaService{}
}

// The initial (s *mediaService) is a method receiver, which means that this function is associated with the mediaService struct.
func (s *mediaService) Upload(
	ctx context.Context,
	request models.UploadMediaRequest) (*models.UploadMediaResult, error) {
	// Starting the implementation of the Upload method. This is where the logic for handling the upload request will be placed.
	cfg, _, err := image.DecodeConfig(request.File)

	if err != nil {
		return nil, err
	}

	_, err = request.File.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return &models.UploadMediaResult{
		OriginalFilename: request.OriginalFilename,
		MimeType:         request.MimeType,
		FileSize:         request.FileSize,
		Width:            cfg.Width,
		Height:           cfg.Height,
	}, nil

}