package services

import (
	"context"
	"image"
	"log"

	"veterans-go-chi-server/internal/imaging"
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

	const maxFileSize int64 = 15 * 1024 *1024 // 15 MB

	if request.FileSize > maxFileSize {
		return nil, ErrFileTooLarge
	}

	mimeType, err := utils.GetMimeType(request.File)
	if err != nil {
		return nil, err
	}

	if !isAllowedMimeType(mimeType) {
		return nil, ErrUnsupportedMediaType
	}


	width, height, err := utils.GetImageDimensions(request.File)
	if err != nil {
		return nil, err
	}
	if !utils.ImageDimensionsAreValid(width, height) {
		return nil, ErrImageDimensionsOutOfRange
	}


	//* Obtain the image from the request file. This consists of a bunch of bytes that represent the image.
	img, _, err := image.Decode(request.File)
	if err != nil {
		return nil, err
	}

	//* Generate the different image variants from the bunch of bytes that represent the image. This will return a slice of GeneratedVariant structs, each containing the name of the variant and the corresponding image data.
	variants, err := imaging.GenerateVariants(img)
	if err != nil {
		return nil, err
	}

	for _, variant := range variants {
	log.Printf(
		"Generated %s: %d bytes",
		variant.Name,
		len(variant.Data),
	)
}

	objectKey := utils.GenerateObjectKey(request.OriginalFilename)

	processed := &models.ProcessedMedia{
		ObjectKey:        objectKey,
		OriginalFilename: request.OriginalFilename,
		MimeType:         mimeType,
		FileSize:         request.FileSize,
		Width:            width,
		Height:           height,
	}

	err = s.repository.Create(ctx, processed)
	if err != nil {
		return nil, err
	}

	return processed, nil

}

func isAllowedMimeType(mimeType string) bool {
	switch mimeType {
		case "image/jpeg", "image/png":
			return true
		default:
			return false
	}
}
