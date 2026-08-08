package services

import (
	"context"
	"log"
	"veterans-go-chi-server/internal/models"
	"veterans-go-chi-server/internal/repositories"
	"veterans-go-chi-server/internal/utils"
)

// This interface is a contract that defines the methods that any media service implementation must provide.
type MediaService interface {
	Upload(
		ctx context.Context,
		request models.UploadMediaRequest,
	) (*models.ProcessedMedia, error)
}

type mediaService struct {
	repository *repositories.MediaRepository
}

// This function is effectively a constructor. It returns a pointer to a mediaService struct, which implements the MediaService interface.
func NewMediaService(repository *repositories.MediaRepository) MediaService {
	return &mediaService{
		repository: repository,
	}
}

// The initial (s *mediaService) is a method receiver, which means that this function is associated with the mediaService struct.
func (s *mediaService) Upload(
	ctx context.Context,
	request models.UploadMediaRequest,
) (*models.ProcessedMedia, error) {
	width, height, err := utils.GetImageDimensions(request.File)
	if err != nil {
		return nil, err
	}
	objectKey := utils.GenerateObjectKeyy(request.OriginalFilename)

	processed := &models.ProcessedMedia{
		ObjectKey:        objectKey,
		OriginalFilename: request.OriginalFilename,
		MimeType:         request.MimeType,
		FileSize:         request.FileSize,
		Width:            width,
		Height:           height,
	}

	log.Printf("BEFORE CREATE: %+v", processed)
	err = s.repository.Create(ctx, processed)
	if err != nil {
		return nil, err
	}
	log.Printf("AFTER CREATE: %+v", processed)

	return processed, nil

}
