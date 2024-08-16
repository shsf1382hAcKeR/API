package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/shsf1382hAcKeR/Canvasify/internal/models"
	"github.com/shsf1382hAcKeR/Canvasify/internal/services"
	"github.com/shsf1382hAcKeR/Canvasify/internal/util"
)

// GetCanvas handles the request to fetch a canvas for a given track URI or URL.
func GetCanvas(w http.ResponseWriter, r *http.Request) {
	// Set the Content-Type header
	w.Header().Set("Content-Type", "application/json")

	var trackURI string

	// Check if track_uri is provided in the query params
	if r.URL.Query().Has("track_uri") {
		trackURI = r.URL.Query().Get("track_uri")
	} else {
		// Check if track_uri is provided in the request body
		var reqBody models.RequestBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Error().Err(err).Msg("❌ Invalid request payload")
			util.WriteError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		trackURI = reqBody.TrackURI
	}

	if trackURI == "" {
		log.Error().Msg("❌ track_uri is required")
		util.WriteError(w, http.StatusBadRequest, "track_uri is required")
		return
	}

	trackID, err := util.ExtractTrackID(trackURI)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		util.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	canvas, err := services.FetchCanvas(trackID)
	if err != nil {
		log.Error().Err(err).Msg("❌ Error fetching canvas")
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := models.SuccessResponse{
		CanvasURL:  canvas.CanvasURL,
		TrackID:    trackID,
		TrackName:  canvas.TrackName,
		ArtistName: canvas.ArtistName,
		PreviewURL: canvas.PreviewURL, // Add this field to the response
	}
	log.Info().Str("ID", trackID).Msg("✔️ Canvas fetched successfully ||")
	util.WriteJSON(w, http.StatusOK, response)
}
