package handlers

import (
	tpl "github.com/Joao-Felisberto/devprivops-dashboard/templates"
	"github.com/labstack/echo"
)

func Hello(c echo.Context) error {
	return tpl.Hello().Render(c.Request().Context(), c.Response())
}
