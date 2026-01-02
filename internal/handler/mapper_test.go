package handler_test

import (
	"testing"
	"time"

	"github.com/serroba/features/internal/flags"
	"github.com/serroba/features/internal/handler"
	"github.com/stretchr/testify/assert"
)

func TestToFlag(t *testing.T) {
	t.Parallel()

	boolVal := true
	body := handler.CreateFlagBody{
		Key:     "test-flag",
		Type:    "bool",
		Enabled: true,
		DefaultValue: handler.ValueBody{
			Kind: "bool",
			Bool: &boolVal,
		},
		Rules: []handler.RuleBody{
			{
				ID: "rule-1",
				Conditions: []handler.ConditionBody{
					{Attr: "plan", Op: "eq", Value: "premium"},
				},
				Value: handler.ValueBody{
					Kind: "bool",
					Bool: &boolVal,
				},
			},
		},
	}

	flag := handler.ToFlag(body)

	assert.Equal(t, "test-flag", flag.Key)
	assert.Equal(t, flags.FlagBool, flag.Type)
	assert.True(t, flag.Enabled)
	assert.Equal(t, flags.FlagBool, flag.DefaultValue.Kind)
	assert.True(t, *flag.DefaultValue.Bool)
	assert.Len(t, flag.Rules, 1)
	assert.Equal(t, "rule-1", flag.Rules[0].ID)
	assert.Len(t, flag.Rules[0].Conditions, 1)
	assert.Equal(t, "plan", flag.Rules[0].Conditions[0].Attr)
	assert.Equal(t, flags.OpEquals, flag.Rules[0].Conditions[0].Op)
}

func TestToFlag_EmptyRules(t *testing.T) {
	t.Parallel()

	body := handler.CreateFlagBody{
		Key:     "simple-flag",
		Type:    "string",
		Enabled: true,
	}

	flag := handler.ToFlag(body)

	assert.Equal(t, "simple-flag", flag.Key)
	assert.Nil(t, flag.Rules)
}

func TestToEvalContext(t *testing.T) {
	t.Parallel()

	body := handler.EvaluateFlagBody{
		TenantID: "tenant-123",
		UserID:   "user-456",
		Attrs: map[string]any{
			"plan":    "premium",
			"country": "US",
		},
	}

	ctx := handler.ToEvalContext(body)

	assert.Equal(t, "tenant-123", ctx.TenantID)
	assert.Equal(t, "user-456", ctx.UserID)
	assert.Equal(t, "premium", ctx.Attrs["plan"])
	assert.Equal(t, "US", ctx.Attrs["country"])
}

func TestToEvalResultBody(t *testing.T) {
	t.Parallel()

	boolVal := true
	now := time.Now()
	result := &flags.EvalResult{
		FlagKey: "my-flag",
		Value: flags.Value{
			Kind: flags.FlagBool,
			Bool: &boolVal,
		},
		Reason:      flags.ReasonRuleMatch,
		RuleID:      "rule-1",
		Version:     5,
		EvaluatedAt: now,
	}

	body := handler.ToEvalResultBody(result)

	assert.Equal(t, "my-flag", body.FlagKey)
	assert.Equal(t, "bool", body.Value.Kind)
	assert.True(t, *body.Value.Bool)
	assert.Equal(t, "rule_match", body.Reason)
	assert.Equal(t, "rule-1", body.RuleID)
	assert.Equal(t, int64(5), body.Version)
	assert.Equal(t, now, body.EvaluatedAt)
}
