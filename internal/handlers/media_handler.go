package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"veterans-go-chi-server/internal/models"
	"veterans-go-chi-server/internal/services"
)

type MediaHandler struct {
	service services.MediaService
}

func NewMediaHandler(service services.MediaService) *MediaHandler {
	return &MediaHandler{
		service: service,
	}
}

func (h *MediaHandler) Upload(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file from the request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	displayOrder, err := getOptionalIntFormValue(r, "display_order")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uploadRequest := models.UploadMediaRequest{
		File:             file,
		OriginalFilename: header.Filename,
		FileSize:         header.Size,
		MediaContext:     getOptionalFormValue(r, "media_context"),
		Season:           getOptionalFormValue(r, "season"),
		Category:         getOptionalFormValue(r, "category"),
		DisplayOrder:     displayOrder,
	}

	processedMedia, err := h.service.Upload(
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
		if errors.Is(err, services.ErrImageDimensionsOutOfRange) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(processedMedia)
}

func getOptionalFormValue(r *http.Request, key string) *string {
    value := r.FormValue(key)
    if value == "" {
        return nil
    }

    return &value
}

func getOptionalIntFormValue(r *http.Request, key string) (*int, error) {
    value := r.FormValue(key)
    if value == "" {
        return nil, nil
    }

    parsed, err := strconv.Atoi(value)
    if err != nil {
        return nil, fmt.Errorf("%s must be an integer", key)
    }

    return &parsed, nil
}
