package handlers

import (
	"encoding/json"
	"io"
	"net/http"
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
		http.Error(w, "Error retrieving the file from the request: " + err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mimeType := http.DetectContentType(buffer)

	uploadRequest := models.UploadMediaRequest {
		File: 				file,
		OriginalFilename: 	header.Filename,
		MimeType: 			mimeType,
		FileSize: 			header.Size,
	}

	result, err := h.service.Upload(
		r.Context(),
		 uploadRequest,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.UploadMediaResponse{
		Filename: result.OriginalFilename,
		Size: result.FileSize,
		MimeType: result.MimeType,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}