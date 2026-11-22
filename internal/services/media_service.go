package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"log"

	"veterans-go-chi-server/internal/constants"
	"veterans-go-chi-server/internal/models"
	"veterans-go-chi-server/internal/repositories"
	"veterans-go-chi-server/internal/storage"
	"veterans-go-chi-server/internal/utils"

	"github.com/jackc/pgx/v5"
)

// This interface is a contract that defines the methods that any media service implementation must provide.
type MediaService interface {
	UploadImage(ctx context.Context, request models.UploadMediaRequest) (*models.ProcessedMedia, error)
	GetImageMetadataById(ctx context.Context, id int) (*models.ImageMetadata, error)
	ListMedia(ctx context.Context, listMediaRequest models.ListMediaRequest) ([]models.MediaItem, error)
	DeleteImageById(ctx context.Context, id int) error
	UpdateImage(ctx context.Context, id int, updateRequest models.UpdateMediaItemRequest) (*models.UpdatedMediaItem, error)
	UpdateImageDisplayPositions(ctx context.Context, updatePositionsRequest models.UpdateDisplayPositionsRequest) error
}

// This struct is the blueprint containing the dependencies required for the media service to function.
type mediaService struct {
	repository    *repositories.MediaRepository
	storage       storage.MediaStorage
	publicBaseURL string
}

// This function is effectively a constructor. It returns a pointer to a mediaService struct, which implements the MediaService interface.
func NewMediaService(
	repository *repositories.MediaRepository,
	storage storage.MediaStorage,
	publicBaseURL string,
) MediaService {
	return &mediaService{
		repository:    repository,
		storage:       storage,
		publicBaseURL: publicBaseURL,
	}
}

const maxImageFileSize int64 = 15 * 1024 * 1024 // 15 MB

var allowedImageMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

// The initial (s *mediaService) is a method receiver, which means that this function is associated with the mediaService struct, and therefore can access its fields and methods.
func (s *mediaService) UploadImage(
	ctx context.Context,
	request models.UploadMediaRequest,
) (processed *models.ProcessedMedia, err error) {

	if err := validateUploadMetadata(request); err != nil {
		return nil, err
	}

	err = s.validateDisplayPosition(ctx, request.DisplayPosition, request.MediaContext, request.Season, request.Category)
	if err != nil {
		return nil, err
	}

	savedKeys := make([]string, 0)

	// rollback: any error return after this point cleans up whatever was already saved to storage
	defer func() {
		if err != nil {
			executeRollback(ctx, s.storage, savedKeys)
		}
	}()

	mimeType, width, height, err := validateImageProperties(request)
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
	variants, err := utils.GenerateVariants(img)
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
		DisplayPosition:  request.DisplayPosition,
		Visible:          request.Visible,
	}

	// Save the processed media information in the database
	if err = s.repository.Create(ctx, processed); err != nil {
		return nil, err
	}

	return processed, nil

}

func (s *mediaService) GetImageMetadataById(ctx context.Context, id int) (*models.ImageMetadata, error) {
	media, err := s.repository.GetImageMetadataById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMediaNotFound
		}
		return nil, err
	}
	return media, nil
}

func (s *mediaService) ListMedia(
	ctx context.Context,
	listMediaRequest models.ListMediaRequest,
) ([]models.MediaItem, error) {
	if err := validateListMediaRequest(listMediaRequest); err != nil {
		return nil, err
	}
	processedMedia, err := s.repository.ListMedia(
		ctx,
		listMediaRequest.MediaContext,
		listMediaRequest.Season,
		listMediaRequest.Category,
	)
	if err != nil {
		return nil, err
	}
	response := make([]models.MediaItem, 0, len(processedMedia))
	for _, media := range processedMedia {
		visible := true
		if media.Visible != nil {
			visible = *media.Visible
		}
		item := models.MediaItem{
			FileDescription: media.FileDescription,
			Width:           media.Width,
			Height:          media.Height,
			Season:          media.Season,
			Category:        media.Category,
			DisplayPosition: media.DisplayPosition,
			Visible:         visible,
		}
		item.URLs.Blur = s.publicBaseURL + "/" + media.BlurKey
		item.URLs.Small = s.publicBaseURL + "/" + media.SmallKey
		item.URLs.Medium = s.publicBaseURL + "/" + media.MediumKey
		item.URLs.Large = s.publicBaseURL + "/" + media.LargeKey
		response = append(response, item)
	}
	return response, nil
}

func (s *mediaService) DeleteImageById(ctx context.Context, id int) error {
	media, err := s.repository.GetImageByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMediaNotFound
		}
		return err
	}
	objectId := media.ObjectKey
	log.Printf("Deleting object with ObjectID: %s", objectId)
	err = s.storage.Delete(ctx, objectId)
	if err != nil {
		return err
	}
	err = s.repository.DeleteImage(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *mediaService) UpdateImageDisplayPositions(
	ctx context.Context,
	updatePositionsRequest models.UpdateDisplayPositionsRequest,
) error {
	// Validate position item list received
	err := validateUpdatePositionSequence(updatePositionsRequest.NewPositions)
	if err != nil {
		return err
	}
	// Invoke repository method to update display positions
	return s.repository.UpdateDisplayPositions(ctx, updatePositionsRequest.NewPositions)
}

func (s *mediaService) UpdateImage(
	ctx context.Context,
	id int,
	updateRequest models.UpdateMediaItemRequest,
) (*models.UpdatedMediaItem, error) {

	// Validar mediaContext y parámetros asociados si se proporciona
	if updateRequest.MediaContext != nil {
		err := validateMediaContextAndParams(*updateRequest.MediaContext, updateRequest.Season, updateRequest.Category)
		if err != nil {
			return nil, err
		}
	}

	return s.repository.UpdateMedia(ctx, id, updateRequest)
}

func validateUpdatePositionSequence(positionItems []models.PositionItem) error {
	if len(positionItems) == 0 {
		return &PositionSequenceError{
			Message:   "position items list cannot be empty",
			ErrorType: "empty",
		}
	}
	seenPositions := make(map[int]bool)

	// Check no duplicate positions
	for _, item := range positionItems {
		if seenPositions[item.Position] {
			return &PositionSequenceError{
				Message:   fmt.Sprintf("duplicate position %d", item.Position),
				ErrorType: "duplicate",
				ErrorPos:  item.Position,
			}
		}
		seenPositions[item.Position] = true
	}

	// Verify sequence is 1, 2, 3, ..., n without gaps
	for i := 1; i <= len(positionItems); i++ {
		if !seenPositions[i] {
			return &PositionSequenceError{
				Message:    fmt.Sprintf("position sequence has gaps, missing position %d", i),
				ErrorType:  "gap",
				MissingPos: i,
			}
		}
	}

	return nil
}

func validateListMediaRequest(request models.ListMediaRequest) error {

	if !isValidMediaContext(request.MediaContext) {
		return ErrInvalidMediaContext
	}

	switch request.MediaContext {
	case string(constants.MediaContextSeason):
		// Season REQUIRES season, category is OPTIONAL
		if request.Season == nil {
			return ErrSeasonRequiredForGivenContext
		}
		if !isValidSeason(*request.Season) {
			return ErrInvalidSeason
		}
		// If category is provided, it must be valid
		if request.Category != nil && !isValidCategory(*request.Category) {
			return ErrInvalidCategory
		}

	case string(constants.MediaContextHome), string(constants.MediaContextBoard), string(constants.MediaContextHistoricMember):
		// Home, Board, and Historical do NOT allow season or category
		if request.Season != nil {
			return ErrSeasonNotAllowedForGivenContext
		}
		if request.Category != nil {
			return ErrCategoryNotAllowedForGivenContext
		}
	}

	return nil
}

func (s *mediaService) validateDisplayPosition(
	ctx context.Context,
	displayPos int,
	mediaContext string,
	season *string,
	category *string,
) error {
	usedPositions, err := s.repository.GetUsedDisplayPositions(
		ctx, mediaContext, season, category,
	)
	if err != nil {
		return err
	}

	availablePositions := inferUnusedDisplayPositions(usedPositions)
	nextPosition := nextDisplayPositionAfterMax(usedPositions)
	availablePositions = append(availablePositions, nextPosition)

	allowed := false
	for _, pos := range availablePositions {
		if pos == displayPos {
			allowed = true
			break
		}
	}
	if !allowed {
		return &DisplayPositionError{
			Message:              "display position is either not allowed or already used",
			AvailablePositions:   availablePositions[:len(availablePositions)-1],
			NextPositionAfterMax: nextPosition,
		}
	}
	return nil
}

func validateUploadMetadata(request models.UploadMediaRequest) error {
	if !isValidMediaContext(request.MediaContext) {
		return ErrInvalidMediaContext
	}

	switch request.MediaContext {
	case string(constants.MediaContextSeason):
		err := validateRequestForSeasonContext(request.Season, request.Category)
		if err != nil {
			return err
		}

	case string(constants.MediaContextHome), string(constants.MediaContextBoard), string(constants.MediaContextHistoricMember):
		err := validateRequestForNonSeasonContext(request.Season, request.Category)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateRequestForSeasonContext(season, category *string) error {
	if season == nil {
		return ErrSeasonRequiredForGivenContext
	}
	if !isValidSeason(*season) {
		return ErrInvalidSeason
	}
	if category == nil {
		return ErrCategoryRequiredForGivenContext
	}
	if !isValidCategory(*category) {
		return ErrInvalidCategory
	}
	if *category != string(constants.BasesCategory) && *category != string(constants.ProCategory) && *category != string(constants.VetsCategory) {
		return ErrCategoryNotAllowedForGivenContext
	}
	return nil
}

func validateRequestForNonSeasonContext(season, category *string) error {
	if season != nil {
		return ErrSeasonNotAllowedForGivenContext
	}
	if category != nil {
		return ErrCategoryNotAllowedForGivenContext
	}
	return nil
}

func validateMediaContextAndParams(mediaContext string, season, category *string) error {
	if !isValidMediaContext(mediaContext) {
		return ErrInvalidMediaContext
	}

	switch mediaContext {
	case string(constants.MediaContextSeason):
		err := validateRequestForSeasonContext(season, category)
		if err != nil {
			return err
		}
	case string(constants.MediaContextHome), string(constants.MediaContextBoard), string(constants.MediaContextHistoricMember):
		err := validateRequestForNonSeasonContext(season, category)
		if err != nil {
			return err
		}
	}
	return nil
}

func inferUnusedDisplayPositions(used []int) []int {
	if len(used) == 0 {
		return []int{}
	}

	maxUsed := used[len(used)-1]
	unused := make([]int, 0)
	usedIdx := 0
	for i := 1; i < maxUsed; i++ {
		if usedIdx < len(used) && used[usedIdx] == i {
			usedIdx++
			continue
		}
		unused = append(unused, i)
	}
	return unused
}

func nextDisplayPositionAfterMax(used []int) int {
	maxUsed := 0
	for _, v := range used {
		if v > maxUsed {
			maxUsed = v
		}
	}
	return maxUsed + 1
}

func executeRollback(ctx context.Context, storage storage.MediaStorage, keys []string) {
	for _, key := range keys {
		if err := storage.Delete(ctx, key); err != nil {
			log.Printf("rollback: failed to delete %s: %v", key, err)
		}
	}
}
