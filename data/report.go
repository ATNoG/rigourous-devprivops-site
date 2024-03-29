package data

type Report struct {
	Branch      string        `json:"branch"`
	Commit      string        `json:"commit"`
	Project     string        `json:"project"`
	Regulations []*Regulation `json:"regulations"`
}

type Regulation struct {
	Name               string        `json:"name"`
	ConsistencyResults []*RuleResult `json:"consistency"`
	PolicyResults      []*RuleResult `json:"policies"`
}

type RuleResult struct {
	Name    string                   `json:"name"`
	Results []map[string]interface{} `json:"results"`
}
