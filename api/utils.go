package main

import (
	"log"
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

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL unset")
	}

	return &EnvVars{
		Domain:      domain,
		Port:        port,
		DatabaseURL: databaseURL,
	}
}
