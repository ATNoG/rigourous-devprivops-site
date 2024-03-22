package handlers

import (
	tpl "github.com/Joao-Felisberto/devprivops-dashboard/templates"
	"github.com/labstack/echo"
)

func Hello(c echo.Context) error {
	return tpl.Hello().Render(c.Request().Context(), c.Response())
}

func Page(c echo.Context) error {
	return tpl.Page("Main page").Render(c.Request().Context(), c.Response())
}
