package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/shsf1382hAcKeR/Canvasify/internal/config"
	"github.com/shsf1382hAcKeR/Canvasify/internal/models"
	"github.com/shsf1382hAcKeR/Canvasify/proto"
	protobuf "google.golang.org/protobuf/proto"
)

// FetchAccessToken retrieves an access token using the refresh token.
func FetchAccessToken() (string, error) {
	form := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {config.RefreshToken},
	}

	req, err := http.NewRequest("POST", config.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(config.ClientID, config.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get access token")
	}

	var respData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	return respData["access_token"].(string), nil
}

// FetchCanvasToken retrieves a canvas token for the web player.
func FetchCanvasToken() (string, error) {
	resp, err := http.Get(config.CanvasTokenURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data["accessToken"].(string), nil
}

// FetchTrackInfo retrieves the track information for a given track ID.
func FetchTrackInfo(trackID string, token string) (string, string, string, string, error) {
	req, _ := http.NewRequest("GET", config.TrackInfoURL+trackID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", "", errors.New("failed to fetch track info")
	}

	var trackData struct {
		ID      string `json:"id"`
		URI     string `json:"uri"`
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		PreviewURL string `json:"preview_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&trackData); err != nil {
		return "", "", "", "", err
	}

	artistNames := make([]string, len(trackData.Artists))
	for i, artist := range trackData.Artists {
		artistNames[i] = artist.Name
	}
	artistName := strings.Join(artistNames, ", ")

	return trackData.URI, trackData.Name, artistName, trackData.PreviewURL, nil
}

// FetchCanvas retrieves the canvas data for a given track ID.
func FetchCanvas(trackID string) (models.Canvas, error) {
	token, err := FetchCanvasToken()
	if err != nil {
		return models.Canvas{}, err
	}

	trackURI, trackName, artistName, previewURL, err := FetchTrackInfo(trackID, token)
	if err != nil {
		return models.Canvas{}, err
	}

	req := &proto.CanvasRequest{}
	req.Tracks = append(req.Tracks, &proto.CanvasRequest_Track{TrackUri: trackURI})
	reqBytes, err := protobuf.Marshal(req)
	if err != nil {
		return models.Canvas{}, err
	}

	request, err := http.NewRequest("POST", config.CanvasFetchURL, bytes.NewReader(reqBytes))
	if err != nil {
		return models.Canvas{}, err
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/protobuf")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return models.Canvas{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Canvas{}, errors.New("failed to fetch canvas")
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Canvas{}, err
	}

	var canvasResp proto.CanvasResponse
	if err := protobuf.Unmarshal(respBytes, &canvasResp); err != nil {
		return models.Canvas{}, err
	}

	for _, canvas := range canvasResp.Canvases {
		if canvas.TrackUri == trackURI {
			return models.Canvas{
				CanvasURL:  canvas.CanvasUrl,
				TrackID:    trackID,
				TrackName:  trackName,
				ArtistName: artistName,
				PreviewURL: previewURL,
			}, nil
		}
	}

	return models.Canvas{}, errors.New("canvas not found")
}
