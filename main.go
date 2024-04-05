package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Joao-Felisberto/devprivops-dashboard/data"
	"github.com/Joao-Felisberto/devprivops-dashboard/handlers"
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

	e.GET("/print/:proj", handlers.PrintPage(store))

	e.POST("/report", handlers.PostReport(store))

	e.Logger.Fatal(e.Start("localhost:8080"))

	store.ToFile("db.json")
}
