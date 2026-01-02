package flags

import (
	"strings"
	"time"
)

type FlagKey string

type FlagType string

const (
	FlagBool   FlagType = "bool"
	FlagString FlagType = "string"
	FlagNumber FlagType = "number"
)

type Flag struct {
	Key          FlagKey
	Type         FlagType
	Enabled      bool // global kill switch
	DefaultValue Value
	Rules        []Rule // ordered: first match wins
	UpdatedAt    time.Time
}

func (f Flag) Evaluate(evalCtx EvalContext) EvalResult {
	result := EvalResult{
		FlagKey:     f.Key,
		EvaluatedAt: time.Now(),
	}

	if !f.Enabled {
		result.Value = f.DefaultValue
		result.Reason = ReasonDisabled

		return result
	}

	for _, rule := range f.Rules {
		if rule.Matches(evalCtx) {
			result.Value = rule.Value
			result.Reason = ReasonRuleMatch
			result.RuleID = rule.ID

			return result
		}
	}

	result.Value = f.DefaultValue
	result.Reason = ReasonDefault

	return result
}

type Rule struct {
	ID         string
	Conditions []Condition // AND across conditions
	Value      Value
}

func (r Rule) Matches(evalCtx EvalContext) bool {
	for _, cond := range r.Conditions {
		if !cond.Matches(evalCtx) {
			return false
		}
	}

	return true
}

type ConditionOp string

const (
	OpEquals     ConditionOp = "eq"
	OpNotEquals  ConditionOp = "neq"
	OpIn         ConditionOp = "in"
	OpNotIn      ConditionOp = "not_in"
	OpExists     ConditionOp = "exists"
	OpStartsWith ConditionOp = "starts_with"
)

type Condition struct {
	Attr  string      // e.g. "tenant_id", "user_id", "plan", "country"
	Op    ConditionOp // eq/in/exists/...
	Value any         // string | float64 | bool | []any depending on Op
}

func (c Condition) Matches(evalCtx EvalContext) bool {
	attrValue := evalCtx.GetAttr(c.Attr)

	switch c.Op {
	case OpEquals:
		return attrValue == c.Value
	case OpNotEquals:
		return attrValue != c.Value
	case OpIn:
		return c.containsValue(attrValue)
	case OpNotIn:
		return !c.containsValue(attrValue)
	case OpExists:
		return attrValue != nil
	case OpStartsWith:
		str, ok := attrValue.(string)
		prefix, prefixOk := c.Value.(string)

		return ok && prefixOk && strings.HasPrefix(str, prefix)
	default:
		return false
	}
}

func (c Condition) containsValue(attrValue any) bool {
	slice, ok := c.Value.([]any)
	if !ok {
		return false
	}

	for _, item := range slice {
		if attrValue == item {
			return true
		}
	}

	return false
}

type Value struct {
	Kind   FlagType
	Bool   *bool
	String *string
	Number *float64
}

func BoolValue(v bool) Value {
	return Value{Kind: FlagBool, Bool: &v}
}

func StringValue(v string) Value {
	return Value{Kind: FlagString, String: &v}
}

func NumberValue(v float64) Value {
	return Value{Kind: FlagNumber, Number: &v}
}

type EvalContext struct {
	TenantID string
	UserID   string
	Attrs    map[string]any // arbitrary attributes for rule conditions
}

func (e EvalContext) GetAttr(attr string) any {
	switch attr {
	case "user_id":
		return e.UserID
	case "tenant_id":
		return e.TenantID
	default:
		if e.Attrs != nil {
			return e.Attrs[attr]
		}

		return nil
	}
}

type EvalReason string

const (
	ReasonDisabled  EvalReason = "disabled"
	ReasonRuleMatch EvalReason = "rule_match"
	ReasonDefault   EvalReason = "default"
)

type EvalResult struct {
	FlagKey     FlagKey
	Value       Value
	Reason      EvalReason
	RuleID      string
	EvaluatedAt time.Time
}
