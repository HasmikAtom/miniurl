package main

import (
	redisdb "github.com/HasmikAtom/miniurl/redis"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"net/http"
)

func (app *App) getOriginalUrl(response http.ResponseWriter, request *http.Request) {
	short := mux.Vars(request)["short"]

	r := redisdb.CreateRedisClient(0)
	defer r.Close()

	fullUrl, err := r.Get(redisdb.Ctx, short).Result()
	if err == redis.Nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(`{"Error": "short not found in db"}`))
		return
	} else if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"Error": "cannot connect to db"}`))
		return
	}

	rInr := redisdb.CreateRedisClient(1)
	defer rInr.Close()

	_ = rInr.Incr(redisdb.Ctx, "counter")

	http.Redirect(response, request, fullUrl, http.StatusFound)
}
