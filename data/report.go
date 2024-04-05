package data

type Report struct {
	Branch      string
	Time        int64
	Project     string
	Regulations []*Regulation
	UserStories []*UserStory
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

type UserStory struct {
	UseCase      string
	IsMisuseCase bool
	Requirements []Requirement
}

type Requirement struct {
	Title       string
	Description string
	Results     []map[string]interface{}
}
