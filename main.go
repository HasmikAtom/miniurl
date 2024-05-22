package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

type ApiServer struct {
	Ctx        context.Context
	DotEnvFile string
	Port       string
	//Database   *database.Database
	// logging    *Logging
}

type EnvVars struct {
	DatabaseURL string
	Port        string
}

type UrlShortenRequest struct {
	Url string `json:"url"`
}
type UrlShortenResponse struct {
	ShortUrl string `json:"shortUrl"`
}

type GetOriginalUrlRequest struct {
}

type GetOriginalUrlresponse struct {
}

func main() {
	envFile := "local.env"
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading env variables from '%s'. Err: %s", envFile, err)
	}

	envVars := GetEnvVars()

	http.Handle("/", handlers())
	log.Printf("Server Running on port :" + envVars.Port)
	err = http.ListenAndServe(":"+envVars.Port, nil)
	if err != nil {
		// graceful shutdown
	}

}

func handlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/shorten", shortenUrl).Methods("POST")

	return router
}

func shortenUrl(response http.ResponseWriter, request *http.Request) {
	log.Printf("Request body =>", request.Body)
	decoder := json.NewDecoder(request.Body)
	body := &UrlShortenRequest{}
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding json =>", err)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(body)
	//response.Write(body)
}

func getOriginalUrl(response http.ResponseWriter, request *http.Request) {

}
