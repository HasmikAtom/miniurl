package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

type App struct {
	Ctx              context.Context
	DotEnvFile       string
	EnvVars          EnvVars
	DefaultUrlExpiry time.Duration
	// logging    *Logging
}

type EnvVars struct {
	Port           string
	Domain         string
	RedisAddr      string
	RedisPass      string
	Quota          string
	QuotaTimeLimit time.Duration
}

type UrlShortenRequest struct {
	Url    string        `json:"url"`
	Expiry time.Duration `json:"expiry"`
}
type UrlShortenResponse struct {
	ShortUrl string        `json:"shortUrl"`
	Expiry   time.Duration `json:"expiry"`
}

func init() {
	envFile := "local.env"
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading env variables from '%s'. Err: %s", envFile, err)
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	envVars := GetEnvVars()
	app := &App{
		Ctx:              ctx,
		EnvVars:          *envVars,
		DefaultUrlExpiry: 24,
	}
	http.Handle("/", handlers(app))
	log.Printf("Server Running on port :" + envVars.Port)
	err := http.ListenAndServe(":"+envVars.Port, nil)
	if err != nil {
		// graceful shutdown
	}

}

func handlers(app *App) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/shorten", app.shortenUrl).Methods("POST")
	router.HandleFunc("/{short}", app.getOriginalUrl).Methods("GET")

	return router
}
