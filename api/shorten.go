package main

import (
	"encoding/json"
	redisdb "github.com/HasmikAtom/miniurl/redis"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"time"
)

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
			response.Write([]byte(`{"Error": "Rate Limit exceeded"}`))
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
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode([]byte(`{"Error": "failed to save the url"}`))
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
