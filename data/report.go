package data

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Joao-Felisberto/devprivops-dashboard/util"
)

type Report struct {
	Branch      string        `json:"branch"`
	Time        int64         `json:"time"`
	Config      string        `json:"config"`
	Project     string        `json:"project"`
	Regulations []*Regulation `json:"policies"`
	UserStories []*UserStory  `json:"user stories"`
	ExtraData   []*ExtraData  `json:"extra data"`
	AttackTrees []*AttackTree `json:"attack trees"`
}

func (r *Report) GetId() string {
	return fmt.Sprintf("%s-%d", strings.ToLower(r.Branch), r.Time)
}

type Regulation struct {
	Name               string        `json:"name"`
	ConsistencyResults []*RuleResult `json:"consistency results"`
	PolicyResults      []*RuleResult `json:"policy results"`
}

func (r *Regulation) UnmarshalJSON(data []byte) error {
	var tmp struct {
		A string        `json:"name"`
		B []*RuleResult `json:"consistency results"`
		C []*RuleResult `json:"policy results"`
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	if !(len(tmp.B) == 0 && len(tmp.C) == 0) {
		*r = Regulation{
			tmp.A,
			tmp.B,
			tmp.C,
		}
	} else {

		var fromReport struct {
			Name    string        `json:"name"`
			Results []*RuleResult `json:"results"`
		}

		err := json.Unmarshal(data, &fromReport)
		if err != nil {
			return err
		}

		*r = Regulation{
			Name:               fromReport.Name,
			ConsistencyResults: util.Filter(fromReport.Results, func(r *RuleResult) bool { return r.IsConsistency }),
			PolicyResults:      util.Filter(fromReport.Results, func(r *RuleResult) bool { return !r.IsConsistency }),
		}
		fmt.Printf("++ %+v\n", r.ConsistencyResults)
		fmt.Printf("-- %+v\n", r.PolicyResults)
	}
	return nil
}

type RuleResult struct {
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	MappingMessage string                   `json:"mapping message"`
	IsConsistency  bool                     `json:"is consistency"`
	Results        []map[string]interface{} `json:"violations"`
}

type UserStory struct {
	UseCase      string        `json:"use case"`
	IsMisuseCase bool          `json:"is misuse case"`
	Requirements []Requirement `json:"requirements"`
}

type Requirement struct {
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Results     []map[string]interface{} `json:"results"`
}

type ExtraData struct {
	Location    string                   `json:"location"`
	Heading     string                   `json:"heading"`
	Description string                   `json:"description"`
	DataRowLine string                   `json:"data row line"`
	Results     []map[string]interface{} `json:"results"`
}

// Represents the execution status of a tree node, either before or after the execution of its associated query
type ExecutionStatus int

const (
	NOT_EXECUTED ExecutionStatus = iota // The node has not yet been executed
	NOT_POSSIBLE                        // The node's condition is deemed not possible
	POSSIBLE                            // The node's condition is deemed possible
	ERROR                               // There was an error when executing the node
)

// Represents a node in the attack tree.
//
// A node is composed of a query, which is its condition, the child nodes and some metadata.
// A node is only evaluated if at least one of its pre-conditions (its children) is possible, or has no children.
type AttackNode struct {
	Description     string                    `json:"description"`      // Brief textual description of the node's condition
	Query           string                    `json:"query"`            // Path to the query that encodes the condition
	Children        []*AttackNode             `json:"children"`         // The node's pre-conditions
	ExecutionStatus ExecutionStatus           `json:"execution status"` // The current execution stats of the node, may change when the tree is executed
	ExecutionResult *[]map[string]interface{} `json:"execution result"` // The result of running the query, if it was run, else nil
}

// Represents the whole attack/harm tree.
//
// Is represented by a singular root node.
// When the root node's condition is possible, the attack/harm is deemed present in the system.
type AttackTree struct {
	Root AttackNode `json:"root"` // The root node of the attack tree
}
