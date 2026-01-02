package handler

import (
	"time"

	"github.com/serroba/features/internal/flags"
)

// Request/Response models for Create Flag

type CreateFlagRequest struct {
	Body CreateFlagBody
}

type CreateFlagBody struct {
	Key          string     `json:"key"                maxLength:"128" minLength:"1" pattern:"^[a-z][a-z0-9-]*$"`
	Type         string     `enum:"bool,string,number" json:"type"`
	Enabled      bool       `json:"enabled"`
	DefaultValue ValueBody  `json:"defaultValue"`
	Rules        []RuleBody `json:"rules,omitempty"`
}

type RuleBody struct {
	ID         string          `json:"id"         maxLength:"64" minLength:"1"`
	Conditions []ConditionBody `json:"conditions" minItems:"1"`
	Value      ValueBody       `json:"value"`
}

type ConditionBody struct {
	Attr  string `json:"attr"                                maxLength:"64" minLength:"1"`
	Op    string `enum:"eq,neq,in,not_in,exists,starts_with" json:"op"`
	Value any    `json:"value"`
}

type ValueBody struct {
	Kind   string   `enum:"bool,string,number" json:"kind"`
	Bool   *bool    `json:"bool,omitempty"`
	String *string  `json:"string,omitempty"`
	Number *float64 `json:"number,omitempty"`
}

type CreateFlagResponse struct {
	Body CreateFlagResponseBody
}

type CreateFlagResponseBody struct {
	Key       flags.FlagKey `json:"key"`
	CreatedAt time.Time     `json:"createdAt"`
}

// Request/Response models for Evaluate Flag

type EvaluateFlagRequest struct {
	Key  string `maxLength:"128" minLength:"1" path:"key" pattern:"^[a-z][a-z0-9-]*$"`
	Body EvaluateFlagBody
}

type EvaluateFlagBody struct {
	TenantID string         `json:"tenantId,omitempty" maxLength:"128"`
	UserID   string         `json:"userId,omitempty"   maxLength:"128"`
	Attrs    map[string]any `json:"attrs,omitempty"`
}

type EvaluateFlagResponse struct {
	Body EvalResultBody
}

type EvalResultBody struct {
	FlagKey     flags.FlagKey `json:"flagKey"`
	Value       ValueBody     `json:"value"`
	Reason      string        `enum:"disabled,rule_match,default" json:"reason"`
	RuleID      string        `json:"ruleId,omitempty"`
	EvaluatedAt time.Time     `json:"evaluatedAt"`
}
