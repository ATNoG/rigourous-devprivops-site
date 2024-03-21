package main

import (
	"html/template"
	"io"

	"github.com/Joao-Felisberto/devprivops-dashboard/handlers"
	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Static("/static", "/static")
	e.GET("/", handlers.Hello)

	e.Logger.Fatal(e.Start("localhost:8080"))
}
