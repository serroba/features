package flags

import "time"

type FlagType string

const (
	FlagBool   FlagType = "bool"
	FlagString FlagType = "string"
	FlagNumber FlagType = "number"
)

type Flag struct {
	Key          string    `json:"key"`
	Type         FlagType  `json:"type"`
	Enabled      bool      `json:"enabled"` // global kill switch
	DefaultValue Value     `json:"default_value"`
	Rules        []Rule    `json:"rules"` // ordered: first match wins
	Version      int64     `json:"version"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Rule struct {
	ID         string      `json:"id"`
	Conditions []Condition `json:"conditions"` // AND across conditions
	Value      Value       `json:"value"`
	Rollout    *Rollout    `json:"rollout,omitempty"` // optional percentage rollout gate
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
	Attr  string      `json:"attr"`  // e.g. "tenant_id", "user_id", "plan", "country"
	Op    ConditionOp `json:"op"`    // eq/in/exists/...
	Value any         `json:"value"` // string | float64 | bool | []any depending on Op
}

type Rollout struct {
	Percentage int    `json:"percentage"` // 0..100
	Salt       string `json:"salt"`       // changes bucket assignment when rotated
	Subject    string `json:"subject"`    // "user" or "tenant" (who gets bucketed)
}

type OverrideScope string

const (
	ScopeUser   OverrideScope = "user"
	ScopeTenant OverrideScope = "tenant"
)

type Override struct {
	FlagKey   string        `json:"flag_key"`
	Scope     OverrideScope `json:"scope"`
	SubjectID string        `json:"subject_id"` // user_id or tenant_id
	Value     Value         `json:"value"`
	Reason    string        `json:"reason,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	CreatedBy string        `json:"created_by,omitempty"`
	ExpiresAt *time.Time    `json:"expires_at,omitempty"` // optional safety valve
}

type Value struct {
	Kind   FlagType `json:"kind"`
	Bool   *bool    `json:"bool,omitempty"`
	String *string  `json:"string,omitempty"`
	Number *float64 `json:"number,omitempty"`
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
	TenantID string         `json:"tenant_id"`
	UserID   string         `json:"user_id"`
	Attrs    map[string]any `json:"attrs,omitempty"` // arbitrary attributes for rule conditions
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
	FlagKey     string     `json:"flag_key"`
	Value       Value      `json:"value"`
	Reason      EvalReason `json:"reason"`
	RuleID      string     `json:"rule_id,omitempty"`
	Version     int64      `json:"version"`
	EvaluatedAt time.Time  `json:"evaluated_at"`
}
