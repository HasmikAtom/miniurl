package main

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetEnvVars() *EnvVars {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "localhost:" + port
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPass := os.Getenv("REDIS_PASS")

	quota := os.Getenv("QUOTA")
	if quota == "" {
		quota = "30"
	}

	timeLimit := os.Getenv("QUOTA_TIME_LIMIT")
	if timeLimit == "" {
		timeLimit = "60"
	}

	minutes, err := strconv.Atoi(timeLimit)
	if err != nil {
		minutes = 60
	}
	quotaTimeLimit := time.Duration(minutes) * time.Minute

	return &EnvVars{
		Domain:         domain,
		Port:           port,
		RedisAddr:      redisAddr,
		RedisPass:      redisPass,
		Quota:          quota,
		QuotaTimeLimit: quotaTimeLimit,
	}
}

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}

func GetIpAddress(request *http.Request) string {
	ipAddress := request.RemoteAddr

	// load balancer/proxy case
	forwardedFor := request.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		ipAddress = forwardedFor
	}

	return ipAddress
}
