package flags

import "time"

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

type Rule struct {
	ID         string
	Conditions []Condition // AND across conditions
	Value      Value
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
