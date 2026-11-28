package services

import (
	"context"
	"io"
	"strings"

	documentmodel "veterans-go-chi-server/internal/models/document"
	"veterans-go-chi-server/internal/repositories"
	"veterans-go-chi-server/internal/storage"
	"veterans-go-chi-server/internal/utils"
)

// MaxDocumentFileSize also bounds the multipart form parsing size in the handler.
const MaxDocumentFileSize int64 = 100 * 1024 * 1024 // 100 MB

var allowedDocumentMimeTypes = map[string]bool{
	"application/pdf": true,
}

type DocumentService interface {
	Upload(
		ctx context.Context,
		request documentmodel.UploadDocumentRequest,
	) (*documentmodel.DocumentMetadata, error)
	List(
		ctx context.Context,
	) ([]documentmodel.DocumentResponse, error)
}

type documentService struct {
	repository    *repositories.DocumentRepository
	storage       storage.MediaStorage
	publicBaseURL string
}

func NewDocumentService(
	repository *repositories.DocumentRepository,
	storage storage.MediaStorage,
	publicBaseURL string,
) DocumentService {
	return &documentService{
		repository:    repository,
		storage:       storage,
		publicBaseURL: publicBaseURL,
	}
}

func (s *documentService) Upload(
	ctx context.Context,
	request documentmodel.UploadDocumentRequest,
) (document *documentmodel.DocumentMetadata, err error) {

	if strings.TrimSpace(request.Title) == "" {
		return nil, ErrDocumentTitleRequired
	}

	mimeType, err := performDocumentValidations(request)
	if err != nil {
		return nil, err
	}

	objectKey := utils.GenerateDocumentObjectKey(request.OriginalFilename)

	// rollback: any error return after this point cleans up whatever was already saved to storage
	saved := false
	defer func() {
		if err != nil && saved {
			executeRollback(ctx, s.storage, []string{objectKey})
		}
	}()

	if _, err = request.File.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(request.File)
	if err != nil {
		return nil, err
	}

	if err = s.storage.Save(ctx, objectKey, data, mimeType); err != nil {
		return nil, err
	}
	saved = true

	document = &documentmodel.DocumentMetadata{
		ObjectKey:        objectKey,
		OriginalFilename: request.OriginalFilename,
		MimeType:         mimeType,
		FileSize:         request.FileSize,
		Title:            request.Title,
		Description:      request.FileDescription,
	}

	if err = s.repository.Create(ctx, document); err != nil {
		return nil, err
	}

	return document, nil
}

func (s *documentService) List(
	ctx context.Context,
) ([]documentmodel.DocumentResponse, error) {
	documents, err := s.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	response := make(
		[]documentmodel.DocumentResponse, 0, len(documents))
	for _, document := range documents {
		response = append(response, documentmodel.DocumentResponse{
			Title:       document.Title,
			Description: document.Description,
			URL:         s.publicBaseURL + "/" + document.ObjectKey,
			MimeType:    document.MimeType,
			FileSize:    document.FileSize,
		})
	}

	return response, nil
}

func performDocumentValidations(request documentmodel.UploadDocumentRequest) (string, error) {
	if err := validateFileSize(request.FileSize, MaxDocumentFileSize); err != nil {
		return "", err
	}

	mimeType, err := utils.GetMimeType(request.File)
	if err != nil {
		return "", err
	}

	if err := validateMimeType(mimeType, allowedDocumentMimeTypes); err != nil {
		return "", err
	}

	return mimeType, nil
}
