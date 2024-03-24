package handlers

// https://templ.guide/syntax-and-usage/context/

import (
	"github.com/Joao-Felisberto/devprivops-dashboard/data"
	tpl "github.com/Joao-Felisberto/devprivops-dashboard/templates"
	"github.com/labstack/echo"
)

func Hello(c echo.Context) error {
	return tpl.Hello().Render(c.Request().Context(), c.Response())
}

func ProjectsPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		return tpl.Page[*data.Report]("Projects Page", tpl.ProjectsPage, store.Data...).Render(ctx.Request().Context(), ctx.Response())
	}
}
