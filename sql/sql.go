package sql

import (
	"errors"
	"fmt"
	"strings"
)

const (
	RuleSetFmt = "(%s)"
	RuleFmt    = "%q %s %v"
)

var (
	InvalidNode = errors.New("invalid Node")
	InvalidRule = errors.New("invalid Rule")
)

var operators = map[string]string{
	"equal":            "=",
	"not_equal":        "!=",
	"in":               "IN",
	"less":             "<",
	"less_or_equal":    "<=",
	"greater":          ">",
	"greater_or_equal": ">=",
}

var conditions = map[string]string{
	"AND": " AND ",
	"OR":  " OR ",
}

var validators = map[string]func(string, interface{}) bool{}

type Node struct {
	*RuleSet
	*Rule
}

type RuleSet struct {
	Condition string  `json:"condition,omitempty"`
	Rules     []*Node `json:"rules,omitempty"`
}

type Rule struct {
	Id       string      `json:"id,omitempty"`
	Field    string      `json:"field,omitempty"`
	Type     string      `json:"type,omitempty"`
	Input    string      `json:"input,omitempty"`
	Operator string      `json:"operator,omitempty"`
	Value    interface{} `json:"value,omitempty"`
}

// SQL format Node to SQL string
func (n *Node) SQL() (string, error) {
	if n.RuleSet != nil {
		return n.RuleSet.SQL()
	} else if n.Rule != nil {
		return n.Rule.SQL()
	}

	return "", InvalidNode
}

// SQL format RuleSet to SQL string
func (r *RuleSet) SQL() (string, error) {
	values := make([]string, 0, len(r.Rules))
	for _, r := range r.Rules {
		value, err := r.SQL()
		if err != nil {
			return "", err
		}
		values = append(values, value)
	}

	return fmt.Sprintf(
		RuleSetFmt,
		strings.Join(values, conditions[r.Condition]),
	), nil
}

// SQL format Rule to SQL string
func (r *Rule) SQL() (string, error) {
	if op, ok := operators[r.Operator]; ok && Validate(r.Field, r.Value) {
		return fmt.Sprintf(RuleFmt, r.Field, op, r.ValueSQL()), nil
	}
	return "", InvalidRule
}

// ValueSQL returns SQL formatted value of Rule
func (r *Rule) ValueSQL() interface{} {
	switch r.Value.(type) {
	case string:
		return fmt.Sprintf("%q", r.Value)
	case nil:
		return "null"
	}
	return r.Value
}

func Validate(field string, value interface{}) bool {
	if validator, ok := validators[field]; ok {
		return validator(field, value)
	}
	return true
}
