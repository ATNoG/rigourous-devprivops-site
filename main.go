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
				Branch:  "a",
				Commit:  "a",
				Project: "a",
				Regulations: []*data.Regulation{
					{
						Name: "Reg 1",
						ConsistencyResults: []*data.RuleResult{
							{
								Name: "Con 1",
								Results: []map[string]interface{}{
									{
										"a": 1,
										"b": 2,
									},
								},
							},
							{
								Name: "Con 2",
								Results: []map[string]interface{}{
									{
										"c": 3,
									},
								},
							},
						},
						PolicyResults: []*data.RuleResult{
							{
								Name: "Pol 1",
								Results: []map[string]interface{}{
									{
										"a": 'a',
										"b": 'b',
										"c": 'c',
									},
								},
							},
						},
					},
					{
						Name: "Reg 2",
						ConsistencyResults: []*data.RuleResult{
							{
								Name: "Con 2",
								Results: []map[string]interface{}{
									{
										"Something": 1,
										"Another":   999,
									},
								},
							},
						},
					},
				},
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
	e.GET("/:proj", handlers.RegulationsPage(&store))
	e.GET("/:proj/:reg", handlers.PoliciesPage(&store))

	e.Logger.Fatal(e.Start("localhost:8080"))
}
