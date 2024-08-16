package models

// RequestBody represents the JSON request body for the GetCanvas endpoint.
type RequestBody struct {
	TrackURI string `json:"track_uri" form:"track_uri"`
}

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

// SuccessResponse represents the structure of a successful response.
type SuccessResponse struct {
	CanvasURL  string `json:"canvas_url"`
	TrackID    string `json:"track_id"`
	TrackName  string `json:"track_name"`
	ArtistName string `json:"artist_name"`
	PreviewURL string `json:"preview_url"`
}

// Canvas represents the data structure for a canvas.
type Canvas struct {
	CanvasURL  string `json:"canvas_url"`
	TrackID    string `json:"track_id"`
	TrackName  string `json:"track_name"`
	ArtistName string `json:"artist_name"`
	PreviewURL string `json:"preview_url"`
}
