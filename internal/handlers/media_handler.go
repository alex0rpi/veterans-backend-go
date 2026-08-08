package handlers

import (
	"encoding/json"
	"net/http"
	"veterans-go-chi-server/internal/models"
	"veterans-go-chi-server/internal/services"
	"veterans-go-chi-server/internal/utils"
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

	mimeType, err := utils.GetMimeType(file)
	if err != nil {
		http.Error(
			w,
			"Error detecting the file mime type: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	uploadRequest := models.UploadMediaRequest{
		File:             file,
		OriginalFilename: header.Filename,
		MimeType:         mimeType,
		FileSize:         header.Size,
	}

	// * Invoque the Upload method of the MediaService to handle the upload request
	processedMedia, err := h.service.Upload(
		r.Context(),
		uploadRequest,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(processedMedia)
}
