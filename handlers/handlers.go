package handlers

// https://templ.guide/syntax-and-usage/context/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func HomePage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		return tpl.PageEmpty(
			"Login",
			func() templ.Component { return tpl.LoginForm() },
		).Render(ctx.Request().Context(), ctx.Response())
	}
}

func Login(store *data.Store) func(ctx echo.Context) error {
	return func(c echo.Context) error {
		levelCookie := new(http.Cookie)
		levelCookie.Name = "level"
		level := c.FormValue("level")
		levelCookie.Value = level
		levelCookie.SameSite = http.SameSiteStrictMode
		c.SetCookie(levelCookie)

		groupCookie := new(http.Cookie)
		groupCookie.Name = "group"
		group := c.FormValue("group")
		groupCookie.Value = group
		groupCookie.SameSite = http.SameSiteStrictMode
		c.SetCookie(groupCookie)

		/*
			levelInt, err := strconv.Atoi(level)
			if err != nil {
				return err
			}
			store.SetAccess(levelInt, group)
		*/

		return tpl.PageEmpty(
			"Login",
			func() templ.Component { return tpl.LoginForm() },
		).Render(c.Request().Context(), c.Response())
	}
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

		lvlCookie, err := ctx.Cookie("level")
		if err != nil {
			return err
		}
		lvl, err := strconv.Atoi(lvlCookie.Value)
		if err != nil {
			return err
		}

		groupCookie, err := ctx.Cookie("group")
		if err != nil {
			return err
		}
		group := groupCookie.Value

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

		filtered := data.Regulation{
			Name: regulation.Name,
			ConsistencyResults: util.Filter(regulation.ConsistencyResults, func(r *data.RuleResult) bool {
				return r.ClearenceLvl <= lvl && util.Contains(r.Groups, group)
			}),
			PolicyResults: util.Filter(regulation.PolicyResults, func(r *data.RuleResult) bool {
				return r.ClearenceLvl <= lvl && util.Contains(r.Groups, group)
			}),
		}

		return tpl.PageSingle[*data.Regulation]("Report Page", func(r *data.Regulation) templ.Component { return tpl.PoliciesPage(project, config, repId, r) }, &filtered).Render(ctx.Request().Context(), ctx.Response())
	}
}

func UserStoriesPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		config := ctx.Param("cfg")
		repId := ctx.Param("repId")

		lvlCookie, err := ctx.Cookie("level")
		if err != nil {
			return err
		}
		lvl, err := strconv.Atoi(lvlCookie.Value)
		if err != nil {
			return err
		}

		groupCookie, err := ctx.Cookie("group")
		if err != nil {
			return err
		}
		group := groupCookie.Value

		report := util.Filter(store.Data, func(r *data.Report) bool {
			return r.Project == project &&
				r.Config == config &&
				r.GetId() == repId
		})[0]

		filtered := util.Filter(report.UserStories, func(us *data.UserStory) bool {
			return us.ClearenceLvl <= lvl && util.Contains(us.Groups, group)
		})

		return tpl.Page[*data.UserStory](
			"User Stories Page",
			func(userStories ...*data.UserStory) templ.Component {
				return tpl.RequirementsPage(project, config, repId, userStories...)
			},
			filtered...).Render(ctx.Request().Context(), ctx.Response())
	}
}

func AttackTreesPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		config := ctx.Param("cfg")
		repId := ctx.Param("repId")

		lvlCookie, err := ctx.Cookie("level")
		if err != nil {
			return err
		}
		lvl, err := strconv.Atoi(lvlCookie.Value)
		if err != nil {
			return err
		}

		groupCookie, err := ctx.Cookie("group")
		if err != nil {
			return err
		}
		group := groupCookie.Value

		report := util.Filter(store.Data, func(r *data.Report) bool {
			return r.Project == project &&
				r.Config == config &&
				r.GetId() == repId
		})[0]

		filtered := util.Filter(report.AttackTrees, func(at *data.AttackTree) bool {
			return at.Root.ClearenceLvl <= lvl && util.Contains(at.Root.Groups, group)
		})

		return tpl.Page[*data.AttackTree](
			"Attack and Harm Tree Page",
			func(trees ...*data.AttackTree) templ.Component {
				return tpl.AttackTreePage(trees...)
			},
			filtered...).Render(ctx.Request().Context(), ctx.Response())
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

		lvlCookie, err := ctx.Cookie("level")
		if err != nil {
			return err
		}
		lvl, err := strconv.Atoi(lvlCookie.Value)
		if err != nil {
			return err
		}

		groupCookie, err := ctx.Cookie("group")
		if err != nil {
			return err
		}
		group := groupCookie.Value

		report := util.Filter(store.Data, func(r *data.Report) bool {
			return r.Project == project &&
				r.Config == config &&
				r.GetId() == repId
		})[0]

		dataList := util.Filter(report.ExtraData, func(d *data.ExtraData) bool {
			return d.Location == dataId && d.ClearenceLvl > lvl && util.Contains(d.Groups, group)
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

		lvlCookie, err := ctx.Cookie("level")
		if err != nil {
			return err
		}
		lvl, err := strconv.Atoi(lvlCookie.Value)
		if err != nil {
			return err
		}

		groupCookie, err := ctx.Cookie("group")
		if err != nil {
			return err
		}
		group := groupCookie.Value

		fmt.Printf("User: %s[%d]\n", group, lvl)

		report := util.Filter(store.Data, func(r *data.Report) bool {
			return r.Project == project &&
				r.Config == config &&
				r.GetId() == repId
		})[0]

		finalReport := data.Report{
			Branch:  report.Branch,
			Time:    report.Time,
			Config:  report.Config,
			Project: report.Project,
			Regulations: util.Map(report.Regulations, func(r *data.Regulation) *data.Regulation {
				return &data.Regulation{
					Name: r.Name,
					ConsistencyResults: util.Filter(r.ConsistencyResults, func(r *data.RuleResult) bool {
						fmt.Printf("con res for lvl %d and %+v\n", r.ClearenceLvl, r.Groups)
						return r.ClearenceLvl <= lvl && util.Contains(r.Groups, group)
					}),
					PolicyResults: util.Filter(r.PolicyResults, func(r *data.RuleResult) bool {
						fmt.Printf("pol res for lvl %d and %+v\n", r.ClearenceLvl, r.Groups)
						return r.ClearenceLvl <= lvl && util.Contains(r.Groups, group)
					}),
				}
			}),
			UserStories: util.Filter(report.UserStories, func(us *data.UserStory) bool {
				return us.ClearenceLvl <= lvl && util.Contains(us.Groups, group)
			}),
			ExtraData: util.Filter(report.ExtraData, func(d *data.ExtraData) bool {
				return d.ClearenceLvl > lvl && util.Contains(d.Groups, group)
			}),
			AttackTrees: util.Filter(report.AttackTrees, func(at *data.AttackTree) bool {
				return at.Root.ClearenceLvl <= lvl && util.Contains(at.Root.Groups, group)
			}),
		}

		return tpl.PageSingle[*data.Report]("Print Page", tpl.PrintPage, &finalReport).Render(ctx.Request().Context(), ctx.Response())
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
