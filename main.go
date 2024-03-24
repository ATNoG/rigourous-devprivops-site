package main

import (
	"github.com/Joao-Felisberto/devprivops-dashboard/data"
	"github.com/Joao-Felisberto/devprivops-dashboard/handlers"
	"github.com/labstack/echo"
)

/*
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
*/

func main() {
	e := echo.New()

	store := data.Store{
		Data: []*data.Report{
			{
				Branch:      "a",
				Commit:      "a",
				Project:     "a",
				Regulations: []*data.Regulation{},
			},
			{
				Branch:      "b",
				Commit:      "b",
				Project:     "b",
				Regulations: []*data.Regulation{},
			},
		},
	}

	e.Static("/static", "static")
	e.GET("/", handlers.ProjectsPage(&store))

	e.Logger.Fatal(e.Start("localhost:8080"))
}
