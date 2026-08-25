package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"veterans-go-chi-server/internal/models"
	"veterans-go-chi-server/internal/services"
)

type DocumentHandler struct {
	service services.DocumentService
}

func NewDocumentHandler(service services.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		service: service,
	}
}

func (h *DocumentHandler) Upload(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseMultipartForm(services.MaxDocumentFileSize); err != nil {
		http.Error(w, "Error parsing the upload: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("document")
	if err != nil {
		http.Error(w, "Error retrieving the file from the request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("error closing file:", err)
		}
	}()

	fileDescription, err := getOptionalDescriptionValue(r, "document_description")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uploadRequest := models.UploadDocumentRequest{
		File:             file,
		OriginalFilename: header.Filename,
		FileSize:         header.Size,
		FileDescription:  fileDescription,
		Title:            r.FormValue("document_title"),
	}

	document, err := h.service.Upload(
		r.Context(),
		uploadRequest,
	)

	if err != nil {
		if errors.Is(err, services.ErrUnsupportedMediaType) {
			http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
			return
		}
		if errors.Is(err, services.ErrFileTooLarge) {
			http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
			return
		}
		if errors.Is(err, services.ErrDocumentTitleRequired) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(document); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *DocumentHandler) List(w http.ResponseWriter, r *http.Request) {
	documents, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(documents); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
