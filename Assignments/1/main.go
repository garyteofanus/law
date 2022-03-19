package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gitlab.com/garyteofanus/law-assignment/auth"
	"gitlab.com/garyteofanus/law-assignment/database"
	"gitlab.com/garyteofanus/law-assignment/http"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	e := echo.New()
	// Some config for echo router
	e.HideBanner = true

	redis := database.NewRedis(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), "")

	authRepo := auth.NewRepo(redis)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	server := http.NewServer(e, os.Getenv("PORT"), authHandler)
	server.Run()
}
