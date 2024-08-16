package util

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// WriteError writes an error response to the client.
func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error_code": statusCode,
		"message":    message,
	})
}

// WriteJSON writes a JSON response to the client.
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// ExtractTrackID extracts the track ID from a given URI or URL.
func ExtractTrackID(trackURI string) (string, error) {
	if strings.HasPrefix(trackURI, "spotify:track:") {
		return strings.TrimPrefix(trackURI, "spotify:track:"), nil
	} else if strings.HasPrefix(trackURI, "https://open.spotify.com/track/") {
		parts := strings.Split(strings.TrimPrefix(trackURI, "https://open.spotify.com/track/"), "?")
		if len(parts) > 0 {
			return parts[0], nil
		}
		return "", errors.New("invalid track link")
	}
	return "", errors.New("invalid track URI or link")
}
