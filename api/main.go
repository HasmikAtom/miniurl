package main

import (
	"context"
	"encoding/json"
	redisdb "github.com/HasmikAtom/miniurl/redis"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

type App struct {
	Ctx        context.Context
	DotEnvFile string
	EnvVars    EnvVars
	// logging    *Logging
}

type EnvVars struct {
	Port      string
	Domain    string
	RedisAddr string
	RedisPass string
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
		Ctx:     ctx,
		EnvVars: *envVars,
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
	router.HandleFunc("/api/{short}", app.shortenUrl).Methods("POST")

	return router
}

func (app *App) shortenUrl(response http.ResponseWriter, request *http.Request) {
	log.Printf("Request body =>", request.Body)
	decoder := json.NewDecoder(request.Body)

	body := &UrlShortenRequest{}
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding json =>", err)
	}

	urlId := uuid.New().String()[:8]

	res := &UrlShortenResponse{
		ShortUrl: app.EnvVars.Domain + "/" + urlId,
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusCreated)
	json.NewEncoder(response).Encode(res)
}

func (app *App) getOriginalUrl(response http.ResponseWriter, request *http.Request) {
	short := mux.Vars(request)["short"]

	r := redisdb.CreateRedisClient(0)
	defer r.Close()

	fullUrl, err := r.Get(redisdb.Ctx, short).Result()
	if err == redis.Nil {
		// short not found in db
	} else if err != nil {
		// cannot connect to db
	}

	rInr := redisdb.CreateRedisClient(1)
	defer rInr.Close()

	_ = rInr.Incr(redisdb.Ctx, "counter")

	response.WriteHeader(http.StatusPermanentRedirect)
	json.NewEncoder(response).Encode(fullUrl)
}
