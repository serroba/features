package flags

import "time"

type FlagType string

const (
	FlagBool   FlagType = "bool"
	FlagString FlagType = "string"
	FlagNumber FlagType = "number"
)

type Flag struct {
	Key          string
	Type         FlagType
	Enabled      bool // global kill switch
	DefaultValue Value
	Rules        []Rule // ordered: first match wins
	Version      int64
	UpdatedAt    time.Time
}

type Rule struct {
	ID         string
	Conditions []Condition // AND across conditions
	Value      Value
	Rollout    *Rollout // optional percentage rollout gate
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

type Rollout struct {
	Percentage int    // 0..100
	Salt       string // changes bucket assignment when rotated
	Subject    string // "user" or "tenant" (who gets bucketed)
}

type OverrideScope string

const (
	ScopeUser   OverrideScope = "user"
	ScopeTenant OverrideScope = "tenant"
)

type Override struct {
	FlagKey   string
	Scope     OverrideScope
	SubjectID string // user_id or tenant_id
	Value     Value
	Reason    string
	CreatedAt time.Time
	CreatedBy string
	ExpiresAt *time.Time // optional safety valve
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
	ReasonDisabled       EvalReason = "disabled"
	ReasonUserOverride   EvalReason = "user_override"
	ReasonTenantOverride EvalReason = "tenant_override"
	ReasonRuleMatch      EvalReason = "rule_match"
	ReasonDefault        EvalReason = "default"
)

type EvalResult struct {
	FlagKey     string
	Value       Value
	Reason      EvalReason
	RuleID      string
	Version     int64
	EvaluatedAt time.Time
}
