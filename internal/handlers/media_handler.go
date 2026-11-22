package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"veterans-go-chi-server/internal/models"
	"veterans-go-chi-server/internal/services"

	"github.com/go-chi/chi/v5"
)

type MediaHandler struct {
	service services.MediaService
}

func NewMediaHandler(service services.MediaService) *MediaHandler {
	return &MediaHandler{
		service: service,
	}
}

func (h *MediaHandler) UploadImage(w http.ResponseWriter, r *http.Request) {

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving the file from the request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("error closing file:", err)
		}
	}()

	mediaContext, err := getRequiredFormValue(r, "media_context")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	displayPosition, err := getIntFormValue(r, "display_position")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileDescription, err := getOptionalDescriptionValue(r, "file_description")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	uploadRequest := models.UploadMediaRequest{
		File:             file,
		OriginalFilename: header.Filename,
		FileSize:         header.Size,
		FileDescription:  fileDescription,
		MediaContext:     mediaContext,
		Season:           getOptionalFormValue(r, "season"),
		Category:         getOptionalFormValue(r, "category"),
		DisplayPosition:  displayPosition,
		Visible:          getOptionalBoolFormValue(r, "visible"),
	}

	processedMedia, err := h.service.UploadImage(
		r.Context(),
		uploadRequest,
	)

	if err != nil {
		// Handle DisplayPositionError with detailed information
		if dpErr, matchesError := err.(*services.DisplayPositionError); matchesError {
			response := map[string]interface{}{
				"error":                     dpErr.Message,
				"available_inner_positions": dpErr.AvailablePositions,
				"next_position_after_last":  dpErr.NextPositionAfterMax,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}
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
		if errors.Is(err, services.ErrInvalidMediaContext) ||
			errors.Is(err, services.ErrInvalidCategory) ||
			errors.Is(err, services.ErrInvalidSeason) ||
			errors.Is(err, services.ErrDisplayPositionNotAllowedOrUsed) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(processedMedia); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MediaHandler) GetImageMetadataById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing media ID", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	metadata, err := h.service.GetImageMetadataById(
		r.Context(),
		idInt,
	)
	if err != nil {
		if errors.Is(err, services.ErrMediaNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metadata); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MediaHandler) ListMedia(w http.ResponseWriter, r *http.Request) {

	var season *string
	seasonParam := r.URL.Query().Get("season")
	if seasonParam != "" {
		season = &seasonParam
	}

	var category *string
	categoryParam := r.URL.Query().Get("category")
	if categoryParam != "" {
		category = &categoryParam
	}

	listMediaRequest := models.ListMediaRequest{
		MediaContext: r.URL.Query().Get("media_context"),
		Season:       season,
		Category:     category,
	}

	log.Printf(
		"REQUEST context=%q seasonParam=%q categoryParam=%q",
		listMediaRequest.MediaContext,
		seasonParam,
		categoryParam,
	)

	// Invoke the service to list media items metadata based on the request parameters
	mediaList, err := h.service.ListMedia(
		r.Context(),
		listMediaRequest,
	)

	if err != nil {
		if errors.Is(err, services.ErrInvalidMediaContext) ||
			errors.Is(err, services.ErrInvalidSeason) ||
			errors.Is(err, services.ErrSeasonRequiredForGivenContext) ||
			errors.Is(err, services.ErrSeasonNotAllowedForGivenContext) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.ListMediaResponse{
		MediaItems: mediaList,
		Total:      len(mediaList),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MediaHandler) DeleteImageById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Missing media ID", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteImageById(
		r.Context(),
		idInt,
	)
	if err != nil {
		if errors.Is(err, services.ErrMediaNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MediaHandler) UpdateImageMetadata(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	// Parse JSON body
	var updateRequest models.UpdateMediaItemRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(
			w,
			"Invalid update request body: "+err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	updatedItem, err := h.service.UpdateImageMetadata(
		r.Context(),
		idInt,
		updateRequest,
	)

	if err != nil {
		if errors.Is(err, services.ErrMediaNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updatedItem)
}

func (h *MediaHandler) UpdateImageDisplayPositions(w http.ResponseWriter, r *http.Request) {
	// Parse JSON body
	var updatePositionsRequest models.UpdateDisplayPositionsRequest
	if err := json.NewDecoder(r.Body).Decode(&updatePositionsRequest); err != nil {
		http.Error(
			w,
			"Invalid display positions update request body: "+err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	err := h.service.UpdateImageDisplayPositions(r.Context(), updatePositionsRequest)
	if err != nil {
		// Validar si es error de secuencia (bad request)
		if posSeqErr, ok := err.(*services.PositionSequenceError); ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := map[string]interface{}{
				"error":      posSeqErr.Message,
				"error_type": posSeqErr.ErrorType,
			}
			switch posSeqErr.ErrorType {
			case "gap":
				errorResponse["missing_position"] = posSeqErr.MissingPos
			case "duplicate":
				errorResponse["duplicate_position"] = posSeqErr.ErrorPos
			}
			_ = json.NewEncoder(w).Encode(errorResponse)
			return
		}
		// Si no es error de secuencia, es error de BD
		http.Error(w, "Failed to update display positions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Display positions updated successfully"})
}

// Secondary functions for handling optional form values and query parameters

func getOptionalBoolFormValue(r *http.Request, key string) *bool {
	value := r.FormValue(key)
	defaultValue := true
	if value == "" {
		return &defaultValue
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return &defaultValue
	}

	return &parsed
}

func getOptionalDescriptionValue(r *http.Request, key string) (*string, error) {
	value := r.FormValue(key)
	if value == "" {
		return nil, nil
	}
	if len(value) > 100 {
		return nil, fmt.Errorf("%s must be at most 100 characters", key)
	}

	return &value, nil
}

func getOptionalFormValue(r *http.Request, key string) *string {
	value := r.FormValue(key)
	if value == "" {
		return nil
	}

	return &value
}

func getRequiredFormValue(r *http.Request, key string) (string, error) {
	value := r.FormValue(key)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}

	return value, nil
}

func getIntFormValue(r *http.Request, key string) (int, error) {
	value := r.FormValue(key)
	if value == "" {
		return 0, fmt.Errorf("%s is required", key)
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer", key)
	}

	return parsed, nil
}
