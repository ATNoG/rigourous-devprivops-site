package main

import (
	"fmt"

	"github.com/Joao-Felisberto/devprivops-dashboard/data"
	"github.com/Joao-Felisberto/devprivops-dashboard/handlers"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	store := data.Store{
		Data: []*data.Report{
			{
				Branch:  "brancha",
				Commit:  "commita",
				Project: "projecta",
				Regulations: []*data.Regulation{
					{
						Name: "Reg1",
						ConsistencyResults: []*data.RuleResult{
							{
								Name: "Con1",
								Results: []map[string]interface{}{
									{
										"a": 1,
										"b": 2,
									},
									{
										"a": 10,
										"b": 20,
									},
								},
							},
							{
								Name: "Con2",
								Results: []map[string]interface{}{
									{
										"c": 3,
									},
								},
							},
						},
						PolicyResults: []*data.RuleResult{
							{
								Name: "Pol1",
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
						Name: "Reg2",
						ConsistencyResults: []*data.RuleResult{
							{
								Name: "Con2",
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
				Branch:      "branchb",
				Commit:      "commitb",
				Project:     "projectb",
				Regulations: []*data.Regulation{},
			},
		},
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

	e.GET("/", handlers.ProjectsPage(&store))
	e.GET("/:proj", handlers.RegulationsPage(&store))
	e.GET("/:proj/:reg", handlers.PoliciesPage(&store))

	e.Logger.Fatal(e.Start("localhost:8080"))
}
