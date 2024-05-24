package main

import (
	"context"
	"encoding/json"
	"fmt"
	redisdb "github.com/HasmikAtom/miniurl/redis"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
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

func (app *App) shortenUrl(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	body := &UrlShortenRequest{}
	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding json =>", err)
	}

	ipAddress := GetIpAddress(request)

	r2 := redisdb.CreateRedisClient(1)
	defer r2.Close()
	val, err := r2.Get(redisdb.Ctx, ipAddress).Result()
	if err == redis.Nil {
		_ = r2.Set(redisdb.Ctx, ipAddress, app.EnvVars.Quota, app.EnvVars.QuotaTimeLimit).Err()
	} else {
		val, _ = r2.Get(redisdb.Ctx, ipAddress).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			response.WriteHeader(http.StatusServiceUnavailable)
			response.Write([]byte(`{Error: Rate Limit exceeded}`))
		}
	}

	body.Url = EnforceHTTP(body.Url)

	urlId := uuid.New().String()[:8]

	if body.Expiry == 0 {
		body.Expiry = app.DefaultUrlExpiry
	}

	r := redisdb.CreateRedisClient(0)
	defer r.Close()

	err = r.Set(redisdb.Ctx, urlId, body.Url, body.Expiry*3600*time.Second).Err()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode([]byte("Error saving the url"))
	}

	res := &UrlShortenResponse{
		ShortUrl: app.EnvVars.Domain + "/" + urlId,
		Expiry:   body.Expiry,
	}

	r2.Decr(redisdb.Ctx, ipAddress)

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

	fmt.Println("FULL URL", fullUrl)

	rInr := redisdb.CreateRedisClient(1)
	defer rInr.Close()

	_ = rInr.Incr(redisdb.Ctx, "counter")

	http.Redirect(response, request, fullUrl, http.StatusFound)
}
