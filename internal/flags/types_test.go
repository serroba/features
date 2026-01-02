package flags_test

import (
	"testing"

	"github.com/serroba/features/internal/flags"
	"github.com/stretchr/testify/assert"
)

func TestEvalContext_GetAttr(t *testing.T) {
	t.Parallel()

	t.Run("returns UserID for user_id", func(t *testing.T) {
		t.Parallel()

		ctx := flags.EvalContext{UserID: "user-123"}
		assert.Equal(t, "user-123", ctx.GetAttr("user_id"))
	})

	t.Run("returns TenantID for tenant_id", func(t *testing.T) {
		t.Parallel()

		ctx := flags.EvalContext{TenantID: "tenant-456"}
		assert.Equal(t, "tenant-456", ctx.GetAttr("tenant_id"))
	})

	t.Run("returns custom attr from Attrs map", func(t *testing.T) {
		t.Parallel()

		ctx := flags.EvalContext{Attrs: map[string]any{"plan": "premium"}}
		assert.Equal(t, "premium", ctx.GetAttr("plan"))
	})

	t.Run("returns nil for missing attr", func(t *testing.T) {
		t.Parallel()

		ctx := flags.EvalContext{Attrs: map[string]any{}}
		assert.Nil(t, ctx.GetAttr("missing"))
	})

	t.Run("returns nil when Attrs is nil", func(t *testing.T) {
		t.Parallel()

		ctx := flags.EvalContext{}
		assert.Nil(t, ctx.GetAttr("anything"))
	})
}

func TestCondition_Matches(t *testing.T) {
	t.Parallel()

	ctx := func(attr string, value any) flags.EvalContext {
		return flags.EvalContext{Attrs: map[string]any{attr: value}}
	}

	t.Run("OpEquals", func(t *testing.T) {
		t.Parallel()

		cond := flags.Condition{Attr: "plan", Op: flags.OpEquals, Value: "premium"}
		assert.True(t, cond.Matches(ctx("plan", "premium")))
		assert.False(t, cond.Matches(ctx("plan", "free")))
	})

	t.Run("OpNotEquals", func(t *testing.T) {
		t.Parallel()

		cond := flags.Condition{Attr: "env", Op: flags.OpNotEquals, Value: "production"}
		assert.True(t, cond.Matches(ctx("env", "staging")))
		assert.False(t, cond.Matches(ctx("env", "production")))
	})

	t.Run("OpIn", func(t *testing.T) {
		t.Parallel()

		cond := flags.Condition{Attr: "region", Op: flags.OpIn, Value: []any{"us-east", "us-west"}}
		assert.True(t, cond.Matches(ctx("region", "us-east")))
		assert.False(t, cond.Matches(ctx("region", "eu-west")))
	})

	t.Run("OpIn with invalid list type", func(t *testing.T) {
		t.Parallel()

		cond := flags.Condition{Attr: "plan", Op: flags.OpIn, Value: "not-a-slice"}
		assert.False(t, cond.Matches(ctx("plan", "premium")))
	})

	t.Run("OpNotIn", func(t *testing.T) {
		t.Parallel()

		cond := flags.Condition{Attr: "country", Op: flags.OpNotIn, Value: []any{"CN", "RU"}}
		assert.True(t, cond.Matches(ctx("country", "US")))
		assert.False(t, cond.Matches(ctx("country", "CN")))
	})

	t.Run("OpNotIn with invalid list type", func(t *testing.T) {
		t.Parallel()

		cond := flags.Condition{Attr: "plan", Op: flags.OpNotIn, Value: "not-a-slice"}
		assert.True(t, cond.Matches(ctx("plan", "premium")))
	})

	t.Run("OpExists", func(t *testing.T) {
		t.Parallel()

		cond := flags.Condition{Attr: "beta", Op: flags.OpExists}
		assert.True(t, cond.Matches(ctx("beta", true)))
		assert.False(t, cond.Matches(flags.EvalContext{Attrs: map[string]any{}}))
	})

	t.Run("OpStartsWith", func(t *testing.T) {
		t.Parallel()

		cond := flags.Condition{Attr: "email", Op: flags.OpStartsWith, Value: "admin@"}
		assert.True(t, cond.Matches(ctx("email", "admin@company.com")))
		assert.False(t, cond.Matches(ctx("email", "user@company.com")))
	})

	t.Run("OpStartsWith with non-string types", func(t *testing.T) {
		t.Parallel()

		cond := flags.Condition{Attr: "num", Op: flags.OpStartsWith, Value: "prefix"}
		assert.False(t, cond.Matches(ctx("num", 123)))
	})

	t.Run("unknown operator", func(t *testing.T) {
		t.Parallel()

		cond := flags.Condition{Attr: "x", Op: "unknown", Value: "y"}
		assert.False(t, cond.Matches(ctx("x", "y")))
	})
}

func TestRule_Matches(t *testing.T) {
	t.Parallel()

	t.Run("matches when all conditions match", func(t *testing.T) {
		t.Parallel()

		rule := flags.Rule{
			Conditions: []flags.Condition{
				{Attr: "plan", Op: flags.OpEquals, Value: "premium"},
				{Attr: "country", Op: flags.OpEquals, Value: "US"},
			},
		}
		ctx := flags.EvalContext{Attrs: map[string]any{"plan": "premium", "country": "US"}}
		assert.True(t, rule.Matches(ctx))
	})

	t.Run("fails when any condition fails", func(t *testing.T) {
		t.Parallel()

		rule := flags.Rule{
			Conditions: []flags.Condition{
				{Attr: "plan", Op: flags.OpEquals, Value: "premium"},
				{Attr: "country", Op: flags.OpEquals, Value: "US"},
			},
		}
		ctx := flags.EvalContext{Attrs: map[string]any{"plan": "premium", "country": "UK"}}
		assert.False(t, rule.Matches(ctx))
	})

	t.Run("matches with empty conditions", func(t *testing.T) {
		t.Parallel()

		rule := flags.Rule{Conditions: []flags.Condition{}}
		assert.True(t, rule.Matches(flags.EvalContext{}))
	})
}

func TestFlag_Evaluate(t *testing.T) {
	t.Parallel()

	t.Run("returns disabled when flag is disabled", func(t *testing.T) {
		t.Parallel()

		flag := flags.Flag{
			Key:          "test",
			Enabled:      false,
			DefaultValue: flags.BoolValue(false),
			Rules: []flags.Rule{
				{ID: "rule-1", Conditions: []flags.Condition{}, Value: flags.BoolValue(true)},
			},
		}
		result := flag.Evaluate(flags.EvalContext{})
		assert.Equal(t, flags.ReasonDisabled, result.Reason)
		assert.False(t, *result.Value.Bool)
	})

	t.Run("returns first matching rule", func(t *testing.T) {
		t.Parallel()

		flag := flags.Flag{
			Key:          "test",
			Enabled:      true,
			DefaultValue: flags.StringValue("default"),
			Rules: []flags.Rule{
				{
					ID:         "rule-1",
					Conditions: []flags.Condition{{Attr: "plan", Op: flags.OpEquals, Value: "enterprise"}},
					Value:      flags.StringValue("first"),
				},
				{
					ID:         "rule-2",
					Conditions: []flags.Condition{{Attr: "plan", Op: flags.OpEquals, Value: "premium"}},
					Value:      flags.StringValue("second"),
				},
			},
		}
		result := flag.Evaluate(flags.EvalContext{Attrs: map[string]any{"plan": "premium"}})
		assert.Equal(t, flags.ReasonRuleMatch, result.Reason)
		assert.Equal(t, "rule-2", result.RuleID)
		assert.Equal(t, "second", *result.Value.String)
	})

	t.Run("returns default when no rules match", func(t *testing.T) {
		t.Parallel()

		flag := flags.Flag{
			Key:          "test",
			Enabled:      true,
			DefaultValue: flags.StringValue("default"),
			Rules: []flags.Rule{
				{
					ID:         "rule-1",
					Conditions: []flags.Condition{{Attr: "plan", Op: flags.OpEquals, Value: "enterprise"}},
					Value:      flags.StringValue("enterprise"),
				},
			},
		}
		result := flag.Evaluate(flags.EvalContext{Attrs: map[string]any{"plan": "free"}})
		assert.Equal(t, flags.ReasonDefault, result.Reason)
		assert.Equal(t, "default", *result.Value.String)
	})

	t.Run("sets metadata correctly", func(t *testing.T) {
		t.Parallel()

		flag := flags.Flag{
			Key:          "my-flag",
			Enabled:      true,
			DefaultValue: flags.BoolValue(true),
		}
		result := flag.Evaluate(flags.EvalContext{})
		assert.Equal(t, flags.FlagKey("my-flag"), result.FlagKey)
		assert.False(t, result.EvaluatedAt.IsZero())
	})
}
