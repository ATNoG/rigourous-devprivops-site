package handlers

// https://templ.guide/syntax-and-usage/context/

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/Joao-Felisberto/devprivops-dashboard/data"
	tpl "github.com/Joao-Felisberto/devprivops-dashboard/templates"
	"github.com/Joao-Felisberto/devprivops-dashboard/util"
	"github.com/a-h/templ"
	"github.com/labstack/echo"
)

func Hello(c echo.Context) error {
	return tpl.Hello("").Render(c.Request().Context(), c.Response())
}

func ProjectsPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		return tpl.Page[*data.Report]("Projects Page", tpl.ProjectsPage, store.Data...).Render(ctx.Request().Context(), ctx.Response())
	}
}

func RegulationsPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		config := ctx.Param("cfg")
		repId := ctx.Param("repId")
		report := util.Filter(store.Data, func(r *data.Report) bool {
			return r.Project == project &&
				r.Config == config &&
				r.GetId() == repId
		})[0]

		return tpl.PageSingle[*data.Report]("Report Page", func(*data.Report) templ.Component { return tpl.RegulationsPage(report, false) }, report).Render(ctx.Request().Context(), ctx.Response())
	}
}

func PoliciesPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		config := ctx.Param("cfg")
		repId := ctx.Param("repId")
		regName := ctx.Param("reg")

		report := util.Filter(store.Data, func(r *data.Report) bool {
			return r.Project == project &&
				r.Config == config &&
				r.GetId() == repId
		})[0]

		for _, rep := range store.Data {
			for _, reg := range rep.Regulations {
				for _, a := range reg.ConsistencyResults {
					fmt.Printf("%+v\n", a)
				}
			}
		}

		fmt.Println(util.Map(report.Regulations, func(r *data.Regulation) int { return len(r.ConsistencyResults) + len(r.PolicyResults) }))

		regulation := util.Filter(report.Regulations, func(r *data.Regulation) bool { return r.Name == regName })[0]

		fmt.Printf("%+v\n", regulation)

		return tpl.PageSingle[*data.Regulation]("Report Page", func(r *data.Regulation) templ.Component { return tpl.PoliciesPage(project, config, repId, r) }, regulation).Render(ctx.Request().Context(), ctx.Response())
	}
}

func UserStoriesPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		config := ctx.Param("cfg")
		repId := ctx.Param("repId")

		report := util.Filter(store.Data, func(r *data.Report) bool {
			return r.Project == project &&
				r.Config == config &&
				r.GetId() == repId
		})[0]

		return tpl.Page[*data.UserStory](
			"User Stories Page",
			func(userStories ...*data.UserStory) templ.Component {
				return tpl.RequirementsPage(project, config, repId, userStories...)
			},
			report.UserStories...).Render(ctx.Request().Context(), ctx.Response())
	}
}

func AttackTreesPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		config := ctx.Param("cfg")
		repId := ctx.Param("repId")

		report := util.Filter(store.Data, func(r *data.Report) bool {
			return r.Project == project &&
				r.Config == config &&
				r.GetId() == repId
		})[0]

		return tpl.Page[*data.AttackTree](
			"Attack and Harm Tree Page",
			func(trees ...*data.AttackTree) templ.Component {
				return tpl.AttackTreePage(trees...)
			},
			report.AttackTrees...).Render(ctx.Request().Context(), ctx.Response())
	}
}

func ExtraData(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		config := ctx.Param("cfg")
		repId := ctx.Param("repId")
		dataId := ctx.Param("id")
		headingLvlStr := ctx.QueryParam("headingLevel")
		headingLevel, err := strconv.Atoi(headingLvlStr)
		if err != nil {
			return err
		}

		report := util.Filter(store.Data, func(r *data.Report) bool {
			return r.Project == project &&
				r.Config == config &&
				r.GetId() == repId
		})[0]

		dataList := util.Filter(report.ExtraData, func(d *data.ExtraData) bool {
			return d.Url == dataId
		})

		if len(dataList) != 0 {
			// return tpl.PageSingle[*data.ExtraData]("Extra Data", func(d *data.ExtraData) templ.Component { return tpl.ExtraData(d, headingLevel) }, dataList[0]).Render(ctx.Request().Context(), ctx.Response())
			return tpl.Page[*data.ExtraData](
				"Extra Data",
				func(d ...*data.ExtraData) templ.Component {
					return tpl.ExtraData(headingLevel, d...)
				}, dataList...).Render(ctx.Request().Context(), ctx.Response())
		}

		return nil
	}
}

func PrintPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		config := ctx.Param("cfg")
		repId := ctx.Param("repId")

		report := util.Filter(store.Data, func(r *data.Report) bool {
			return r.Project == project &&
				r.Config == config &&
				r.GetId() == repId
		})[0]

		return tpl.PageSingle[*data.Report]("Print Page", tpl.PrintPage, report).Render(ctx.Request().Context(), ctx.Response())
	}
}

func PostReport(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		// var raw []byte
		// rawEncoded, err := io.ReadAll(ctx.Request().Body)
		raw, err := io.ReadAll(ctx.Request().Body)
		if err != nil {
			return err
		}

		iface := map[string]interface{}{}
		err = json.Unmarshal(raw, &iface)
		if err != nil {
			return err
		}

		m, err := json.MarshalIndent(iface, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(m))

		var report data.Report
		err = json.Unmarshal(raw, &report)
		if err != nil {
			fmt.Println(err)
			return err
		}

		final, err := json.MarshalIndent(report, "", " ")
		if err != nil {
			return err
		}

		store.Data = append(store.Data, &report)

		store.ToFile("db.json")

		return tpl.Hello(string(final)).Render(ctx.Request().Context(), ctx.Response())
	}
}
