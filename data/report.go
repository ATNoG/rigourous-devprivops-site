package data

type Report struct {
	Branch      string
	Time        int64
	Project     string
	Regulations []*Regulation
}

type Regulation struct {
	Name               string
	ConsistencyResults []*RuleResult
	PolicyResults      []*RuleResult
}

type RuleResult struct {
	Name           string
	MappingMessage string
	Results        []map[string]interface{}
}
