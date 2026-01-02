package handler

import (
	"time"
)

// Request/Response models for Create Flag

type CreateFlagRequest struct {
	Body CreateFlagBody
}

type CreateFlagBody struct {
	Key          string     `json:"key"`
	Type         string     `json:"type"`
	Enabled      bool       `json:"enabled"`
	DefaultValue ValueBody  `json:"defaultValue"`
	Rules        []RuleBody `json:"rules,omitempty"`
}

type RuleBody struct {
	ID         string          `json:"id"`
	Conditions []ConditionBody `json:"conditions"`
	Value      ValueBody       `json:"value"`
}

type ConditionBody struct {
	Attr  string `json:"attr"`
	Op    string `json:"op"`
	Value any    `json:"value"`
}

type ValueBody struct {
	Kind   string   `json:"kind"`
	Bool   *bool    `json:"bool,omitempty"`
	String *string  `json:"string,omitempty"`
	Number *float64 `json:"number,omitempty"`
}

type CreateFlagResponse struct {
	Body CreateFlagResponseBody
}

type CreateFlagResponseBody struct {
	Key       string    `json:"key"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
}

// Request/Response models for Evaluate Flag

type EvaluateFlagRequest struct {
	Key  string `path:"key"`
	Body EvaluateFlagBody
}

type EvaluateFlagBody struct {
	TenantID string         `json:"tenantId,omitempty"`
	UserID   string         `json:"userId,omitempty"`
	Attrs    map[string]any `json:"attrs,omitempty"`
}

type EvaluateFlagResponse struct {
	Body EvalResultBody
}

type EvalResultBody struct {
	FlagKey     string    `json:"flagKey"`
	Value       ValueBody `json:"value"`
	Reason      string    `json:"reason"`
	RuleID      string    `json:"ruleId,omitempty"`
	Version     int64     `json:"version"`
	EvaluatedAt time.Time `json:"evaluatedAt"`
}
