package handlers

// https://templ.guide/syntax-and-usage/context/

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/Joao-Felisberto/devprivops-dashboard/data"
	tpl "github.com/Joao-Felisberto/devprivops-dashboard/templates"
	"github.com/Joao-Felisberto/devprivops-dashboard/util"
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

func PostReport(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		// var raw []byte
		raw, err := io.ReadAll(ctx.Request().Body)
		if err != nil {
			return err
		}

		iface := map[string]interface{}{}
		err = json.Unmarshal(raw, &iface)
		if err != nil {
			return err
		}

		branch := iface["branch"].(string)
		commit := iface["commit"].(string)
		project := iface["project"].(string)

		regulationsRaw := iface["policies"].(map[string]interface{})

		regulationsList := []*data.Regulation{}

		// policies->X
		for regName, reg := range regulationsRaw {
			consistencyRes := []*data.RuleResult{}
			policyRes := []*data.RuleResult{}
			regulation := reg.(map[string]interface{})

			// policies->[gdpr]->X
			for polName, pol := range regulation {
				policy := pol.(map[string]interface{})

				var violations []map[string]interface{}
				if policy["violations"] != nil {
					violationsRaw := policy["violations"].([]interface{})
					for _, v := range violationsRaw {
						violations = append(violations, v.(map[string]interface{}))
					}
				}

				if policy["is consistency"].(bool) {
					consistencyRes = append(consistencyRes, &data.RuleResult{
						Name:    polName,
						Results: violations,
					})
				} else {
					policyRes = append(policyRes, &data.RuleResult{
						Name:    polName,
						Results: violations,
					})
				}
			}

			fmt.Printf("%s cons:%d pol:%d\n", regName, len(consistencyRes), len(policyRes))
			regulationsList = append(regulationsList, &data.Regulation{
				Name:               regName,
				ConsistencyResults: consistencyRes,
				PolicyResults:      policyRes,
			})
		}

		report := data.Report{
			Branch:      branch,
			Commit:      commit,
			Project:     project,
			Regulations: regulationsList,
		}
		final, err := json.MarshalIndent(report, "", " ")
		if err != nil {
			return err
		}

		store.Data = append(store.Data, &report)

		return tpl.Hello(string(final)).Render(ctx.Request().Context(), ctx.Response())
	}
}
