package main

import (
	"os"
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

	return &EnvVars{
		Domain:    domain,
		Port:      port,
		RedisAddr: redisAddr,
		RedisPass: redisPass,
	}
}

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}
