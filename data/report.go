package data

type Report struct {
	Branch      string
	Commit      string
	Project     string
	Regulations []*Regulation
}

type Regulation struct {
	Name               string
	ConsistencyResults []*RuleResult
	PolicyResults      []*RuleResult
}

type RuleResult struct {
	Name    string
	Results []map[string]interface{}
}
