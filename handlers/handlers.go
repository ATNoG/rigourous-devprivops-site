package handlers

// https://templ.guide/syntax-and-usage/context/

import (
	"github.com/Joao-Felisberto/devprivops-dashboard/data"
	tpl "github.com/Joao-Felisberto/devprivops-dashboard/templates"
	"github.com/Joao-Felisberto/devprivops-dashboard/util"
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

func RegulationsPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		report := util.Filter(store.Data, func(r *data.Report) bool { return r.Project == project })[0]

		return tpl.PageSingle[*data.Report]("Report Page", tpl.RegulationsPage, report).Render(ctx.Request().Context(), ctx.Response())
	}
}
func PoliciesPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		regName := ctx.Param("reg")
		report := util.Filter(store.Data, func(r *data.Report) bool { return r.Project == project })[0]
		regulation := util.Filter(report.Regulations, func(r *data.Regulation) bool { return r.Name == regName })[0]

		return tpl.PageSingle[*data.Regulation]("Report Page", tpl.PoliciesPage, regulation).Render(ctx.Request().Context(), ctx.Response())
	}
}
