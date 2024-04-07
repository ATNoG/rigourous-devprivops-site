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
		report := util.Filter(store.Data, func(r *data.Report) bool { return r.Project == project })[0]

		return tpl.PageSingle[*data.Report]("Report Page", func(*data.Report) templ.Component { return tpl.RegulationsPage(report, false) }, report).Render(ctx.Request().Context(), ctx.Response())
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

func UserStoriesPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")

		report := util.Filter(store.Data, func(r *data.Report) bool { return r.Project == project })[0]

		return tpl.Page[*data.UserStory]("User Stories Page", tpl.RequirementsPage, report.UserStories...).Render(ctx.Request().Context(), ctx.Response())
	}
}

func ExtraData(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")
		dataId := ctx.Param("id")
		headingLvlStr := ctx.QueryParam("headingLevel")
		headingLevel, err := strconv.Atoi(headingLvlStr)
		if err != nil {
			return err
		}

		report := util.Filter(store.Data, func(r *data.Report) bool { return r.Project == project })[0]

		dataList := util.Filter(report.ExtraData, func(d *data.ExtraData) bool {
			return d.Url == dataId
		})

		if len(dataList) != 0 {
			return tpl.PageSingle[*data.ExtraData]("Extra Data", func(d *data.ExtraData) templ.Component { return tpl.ExtraData(d, headingLevel) }, dataList[0]).Render(ctx.Request().Context(), ctx.Response())
		}

		fmt.Println("Could not find extra data")
		return nil
	}
}

func PrintPage(store *data.Store) func(ctx echo.Context) error {
	return func(ctx echo.Context) error {
		project := ctx.Param("proj")

		report := util.Filter(store.Data, func(r *data.Report) bool { return r.Project == project })[0]

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
			fmt.Printf("2: %s\n", string(raw))
			return err
		}

		m, err := json.MarshalIndent(iface, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(m))

		branch := iface["branch"].(string)
		time, err := strconv.ParseInt(iface["time"].(string), 10, 64)
		if err != nil {
			// TODO: better error handling
			panic(err)
		}
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
						Name:           polName,
						Results:        violations,
						MappingMessage: policy["mapping message"].(string),
					})
				} else {
					policyRes = append(policyRes, &data.RuleResult{
						Name:           polName,
						Results:        violations,
						MappingMessage: policy["mapping message"].(string),
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

		userStoriesRaw := iface["user stories"].(map[string]interface{})
		userStoryList := []*data.UserStory{}

		// user stories->X
		for usName, usRaw := range userStoriesRaw {
			us := usRaw.(map[string]interface{})

			isMisuseCase := us["is misuse case"].(bool)
			requirementsRaw := us["requirements"].([]interface{})

			requirements := util.Map(requirementsRaw, func(r interface{}) data.Requirement {
				req := r.(map[string]interface{})
				ress := util.Map(req["results"].([]interface{}), func(rr interface{}) map[string]interface{} { return rr.(map[string]interface{}) })
				return data.Requirement{
					Title:       req["title"].(string),
					Description: req["description"].(string),
					Results:     ress,
				}
			})

			userStoryList = append(userStoryList, &data.UserStory{
				UseCase:      usName,
				IsMisuseCase: isMisuseCase,
				Requirements: requirements,
			})
		}

		extraDataRaw := iface["extra data"].([]interface{})
		extraDataList := []*data.ExtraData{}

		// extra data->X
		for _, dRaw := range extraDataRaw {
			d := dRaw.(map[string]interface{})
			// resultsRaw := d["results"].([]interface{})
			resultsRaw := []interface{}{}
			if d["results"] != nil {
				resultsRaw = d["results"].([]interface{})
			}
			results := util.Map(resultsRaw, func(r interface{}) map[string]interface{} { return r.(map[string]interface{}) })

			extraDataList = append(extraDataList, &data.ExtraData{
				Url:         d["url"].(string),
				Heading:     d["heading"].(string),
				Description: d["description"].(string),
				DataRowLine: d["data row line"].(string),
				Results:     results,
			})
		}

		report := data.Report{
			Branch:      branch,
			Time:        time,
			Project:     project,
			Regulations: regulationsList,
			UserStories: userStoryList,
			ExtraData:   extraDataList,
		}
		final, err := json.MarshalIndent(report, "", " ")
		if err != nil {
			return err
		}

		store.Data = append(store.Data, &report)
		fmt.Printf("%s\n", string(final))

		return tpl.Hello(string(final)).Render(ctx.Request().Context(), ctx.Response())
	}
}
