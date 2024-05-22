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

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL unset")
	}

	return &EnvVars{
		Port:        port,
		DatabaseURL: databaseURL,
	}
}
