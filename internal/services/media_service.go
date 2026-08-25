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
	ListMedia(
		ctx context.Context,
		listMediaRequest models.ListMediaRequest,
	) ([]models.ListMediaResponse, error)
}

// This struct is the blueprint containing the dependencies required for the media service to function.
type mediaService struct {
	repository *repositories.MediaRepository
	storage    storage.MediaStorage
	publicBaseURL string
}

// This function is effectively a constructor. It returns a pointer to a mediaService struct, which implements the MediaService interface.
func NewMediaService(
	repository *repositories.MediaRepository,
	storage storage.MediaStorage,
	publicBaseURL string,
) MediaService {
	return &mediaService{
		repository: repository,
		storage:    storage,
		publicBaseURL: publicBaseURL,
	}
}

const maxImageFileSize int64 = 15 * 1024 * 1024 // 15 MB

var allowedImageMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

// The initial (s *mediaService) is a method receiver, which means that this function is associated with the mediaService struct, and therefore can access its fields and methods.
func (s *mediaService) Upload(
	ctx context.Context,
	request models.UploadMediaRequest,
) (processed *models.ProcessedMedia, err error) {

	if err := validateUploadMetadata(request); err != nil {
		return nil, err
	}

	savedKeys := make([]string, 0)

	// rollback: any error return after this point cleans up whatever was already saved to storage
	defer func() {
		if err != nil {
			executeRollback(ctx, s.storage, savedKeys)
		}
	}()

	mimeType, width, height, err := performRequestValidations(request)
	if err != nil {
		return nil, err
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

	//Save the original storage.
	if err := s.storage.Save(
		ctx, objectKey, originalData, mimeType,
	); err != nil {
		return nil, err
	}
	savedKeys = append(savedKeys, objectKey)

	//Save the generated variants to the storage
	variantKeys := make(map[string]string, len(variants))

	for _, variant := range variants {
		log.Printf("Saving variant: %s, size: %d bytes", variant.Name, len(variant.Data))
		variantKey := utils.GenerateVariantObjectKey(objectKey, variant.Name)
		if err := s.storage.Save(
			ctx, variantKey, variant.Data, "image/webp",
		); err != nil {
			return nil, err
		}
		savedKeys = append(savedKeys, variantKey)
		variantKeys[variant.Name] = variantKey
	}

	processed = &models.ProcessedMedia{
		ObjectKey:        objectKey,
		OriginalFilename: request.OriginalFilename,
		FileDescription:  request.FileDescription,
		MimeType:         mimeType,
		FileSize:         request.FileSize,
		Width:            width,
		Height:           height,
		BlurKey:          variantKeys[string(constants.VariantBlur)],
		SmallKey:         variantKeys[string(constants.VariantSmall)],
		MediumKey:        variantKeys[string(constants.VariantMedium)],
		LargeKey:         variantKeys[string(constants.VariantLarge)],
		MediaContext:     request.MediaContext,
		Season:           request.Season,
		Category:         request.Category,
		DisplayOrder:     request.DisplayOrder,
	}

	// Save the processed media information in the database
	if err = s.repository.Create(ctx, processed); err != nil {
		return nil, err
	}

	return processed, nil

}

func (s *mediaService) ListMedia(
	ctx context.Context,
	listMediaRequest models.ListMediaRequest,
) ([]models.ListMediaResponse, error) {
	if err := validateListMediaRequest(listMediaRequest); err != nil {
		return nil, err
	}
	processedMedia, err := s.repository.ListMedia(
		ctx,
		 listMediaRequest.MediaContext,
		  listMediaRequest.Season,
		)
	if err != nil {
		return nil, err
	}
	response := make(
		[]models.ListMediaResponse, 0, len(processedMedia))
	for _, media := range processedMedia {
		response = append(response, models.ListMediaResponse{
			OriginalFilename: media.OriginalFilename,
			FileDescription:  media.FileDescription,
			Width:            media.Width,
			Height:           media.Height,
			BlurMediaURL:     s.publicBaseURL + "/" + media.BlurKey,
			SmallMediaURL:    s.publicBaseURL + "/" + media.SmallKey,
			MediumMediaURL:   s.publicBaseURL + "/" + media.MediumKey,
			LargeMediaURL:    s.publicBaseURL + "/" + media.LargeKey,
			MediaContext:     media.MediaContext,
			Season:           media.Season,
			Category:         media.Category,
			DisplayOrder:     media.DisplayOrder,
		})
	}
	return response, nil
}

func performRequestValidations(request models.UploadMediaRequest) (string, int, int, error) {
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

func executeRollback(ctx context.Context, storage storage.MediaStorage, keys []string) {
	for _, key := range keys {
		if err := storage.Delete(ctx, key); err != nil {
			log.Printf("rollback: failed to delete %s: %v", key, err)
		}
	}
}
