package main

import (
	"log"
	"movieapp/config"
	"movieapp/server"
	"movieapp/tmdb"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Static("/static", "static")

	if err := tmdb.SetUpApi().LoadGenres(); err != nil {
		log.Fatalf("Error loading genres: %v", err)
	}

	db := config.ConnectMongoDB()
	_ = db
	api := tmdb.SetUpApi()
	server.SetUpRoutes(e, api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	e.Logger.Fatal(e.Start(":" + port))

}
