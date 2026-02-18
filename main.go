package main

import (
	"log"
	"movieapp/config"
	"movieapp/server"
	"movieapp/tmdb"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000", // Next.js dev
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
	}))

	e.Static("/static", "static")

	if err := tmdb.SetUpApi().LoadGenres(); err != nil {
		log.Fatalf("Error loading genres: %v", err)
	}

	db := config.ConnectMongoDB()
	_ = db
	api := tmdb.SetUpApi()
	server.SetUpRoutes(e, api, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	e.Logger.Fatal(e.Start(":" + port))

}
