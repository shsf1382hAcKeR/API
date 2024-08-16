package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// "github.com/joho/godotenv"
	"github.com/shsf1382hAcKeR/Canvasify/internal/handlers"
	"github.com/shsf1382hAcKeR/Canvasify/internal/logging"
)

func main() {
	// Define the port flag
	port := flag.String("port", "8080", "Port to run the server on")
	flag.Parse()

	

	// Load the .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	logging.Logger.Fatal().Err(err).Msg("Error loading .env file")
	// }

	// Set up zerolog
	logging.SetupLogger()

	// Create a new router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Define routes
	r.Get("/v1/canvas", handlers.GetCanvas)

	// Start the server
	url := fmt.Sprintf("http://localhost:%s", *port)
	logging.Logger.Info().Msgf("💻 Server is live at %s", logging.ColoredURL(url))
	err := http.ListenAndServe(":"+*port, r)
	if err != nil {
		logging.Logger.Fatal().Err(err).Msg("Error starting server")
	}
}
