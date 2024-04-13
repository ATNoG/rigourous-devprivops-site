package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/Joao-Felisberto/devprivops-dashboard/data"
	"github.com/Joao-Felisberto/devprivops-dashboard/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	store, err := data.FromFile("./db.json")
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
		os.Exit(1)
	}

	e.Static("/static", "static")

	for _, f := range []string{
		"android-chrome-192x192.png",
		"android-chrome-512x512.png",
		"apple-touch-icon.png",
		"favicon-16x16.png",
		"favicon-32x32.png",
		"favicon.ico",
	} {
		e.Static(
			fmt.Sprintf("/%s", f),
			"static",
		)
	}
	e.Static("site.manifest", "/static/site.manifest")

	e.GET("/", handlers.ProjectsPage(store))
	e.GET("/view/:proj", handlers.RegulationsPage(store))
	e.GET("/view/:proj/:reg", handlers.PoliciesPage(store))
	e.GET("/us/:proj", handlers.UserStoriesPage(store))

	e.GET("/data/:proj/:id", handlers.ExtraData(store))

	e.GET("/print/:proj", handlers.PrintPage(store))

	e.POST("/report", handlers.PostReport(store))

	err = godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	host, found := os.LookupEnv("HOST")
	if !found {
		slog.Error("'HOST' key not found in environment")
	}
	port, found := os.LookupEnv("PORT")
	if !found {
		slog.Error("'PORT' key not found in environment")
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", host, port)))

	store.ToFile("db.json")
}
