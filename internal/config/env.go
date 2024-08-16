package config

import (
	"os"
)

// Spotify API credentials
var (
	ClientID     = os.Getenv("SPOTIFY_CLIENT_ID")
	ClientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")
	RefreshToken = os.Getenv("SPOTIFY_REFRESH_TOKEN")
)
