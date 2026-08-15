package services

import (
	"bytes"
	"context"
	"image"
	"io"
	"log"

	"veterans-go-chi-server/internal/constants"
	"veterans-go-chi-server/internal/imaging"
	"veterans-go-chi-server/internal/models"
	"veterans-go-chi-server/internal/repositories"
	"veterans-go-chi-server/internal/storage"
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
	repository 		*repositories.MediaRepository
	storage			storage.MediaStorage
}

// This function is effectively a constructor. It returns a pointer to a mediaService struct, which implements the MediaService interface.
func NewMediaService(
	repository *repositories.MediaRepository,
	storage storage.MediaStorage,
) MediaService {
	return &mediaService{
		repository: repository,
		storage:    storage,
	}
}

// The initial (s *mediaService) is a method receiver, which means that this function is associated with the mediaService struct.
func (s *mediaService) Upload(
	ctx context.Context,
	request models.UploadMediaRequest,
) (*models.ProcessedMedia, error) {

	const maxFileSize int64 = 15 * 1024 * 1024 // 15 MB

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

	objectKey := utils.GenerateObjectKey(request.OriginalFilename)

	// Read the upload once so it can be stored and decoded independently.
	_, err = request.File.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	originalData, err := io.ReadAll(request.File)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(originalData))
	if err != nil {
		return nil, err
	}

	//Generate the different image variants from the bunch of bytes that represent the image. This will return a slice of GeneratedVariant structs, each containing the name of the variant and the corresponding image data.
	variants, err := imaging.GenerateVariants(img)
	if err != nil {
		return nil, err
	}
	//End of variants generation

	//Save the original and generated variants to storage.
	if err := s.storage.Save(ctx, objectKey, originalData); err != nil {
		return nil, err
	}

	//Save the generated variants to the storage
	variantKeys := make(map[string]string, len(variants))

	for _, variant := range variants {
		log.Printf("Saving variant: %s, size: %d bytes", variant.Name, len(variant.Data))
		variantKey := utils.GenerateVariantObjectKey(objectKey, variant.Name)
		if err := s.storage.Save(ctx, variantKey, variant.Data); err != nil {
			return nil, err
		}
		variantKeys[variant.Name] = variantKey
	}

	processed := &models.ProcessedMedia{
		ObjectKey:        objectKey,
		OriginalFilename: request.OriginalFilename,
		MimeType:         mimeType,
		FileSize:         request.FileSize,
		Width:            width,
		Height:           height,
		BlurKey:          variantKeys[constants.VariantBlur],
		SmallKey:         variantKeys[constants.VariantSmall],
		MediumKey:        variantKeys[constants.VariantMedium],
		LargeKey:         variantKeys[constants.VariantLarge],
	}

	// Save the processed media information in the database
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
