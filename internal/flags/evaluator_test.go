package flags_test

import (
	"testing"

	"github.com/serroba/features/internal/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuleMatcher_ReturnsFirstMatch(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
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
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"plan": "premium"},
	})

	require.NotNil(t, result)
	assert.Equal(t, "rule-2", result.ID)
}

func TestRuleMatcher_NoMatch(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "rule-1",
			Conditions: []flags.Condition{{Attr: "plan", Op: flags.OpEquals, Value: "enterprise"}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"plan": "free"},
	})

	assert.Nil(t, result)
}

func TestRuleMatcher_EmptyRules(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	result := matcher([]flags.Rule{}, flags.EvalContext{})

	assert.Nil(t, result)
}

func TestRuleMatcher_OpEquals(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "equals-rule",
			Conditions: []flags.Condition{{Attr: "country", Op: flags.OpEquals, Value: "US"}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"country": "US"},
	})
	require.NotNil(t, result)

	result = matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"country": "UK"},
	})
	assert.Nil(t, result)
}

func TestRuleMatcher_OpNotEquals(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "not-equals-rule",
			Conditions: []flags.Condition{{Attr: "env", Op: flags.OpNotEquals, Value: "production"}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"env": "staging"},
	})
	require.NotNil(t, result)

	result = matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"env": "production"},
	})
	assert.Nil(t, result)
}

func TestRuleMatcher_OpIn(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "in-rule",
			Conditions: []flags.Condition{{Attr: "region", Op: flags.OpIn, Value: []any{"us-east", "us-west"}}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"region": "us-east"},
	})
	require.NotNil(t, result)

	result = matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"region": "eu-west"},
	})
	assert.Nil(t, result)
}

func TestRuleMatcher_OpNotIn(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "not-in-rule",
			Conditions: []flags.Condition{{Attr: "country", Op: flags.OpNotIn, Value: []any{"CN", "RU"}}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"country": "US"},
	})
	require.NotNil(t, result)

	result = matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"country": "CN"},
	})
	assert.Nil(t, result)
}

func TestRuleMatcher_OpExists(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "exists-rule",
			Conditions: []flags.Condition{{Attr: "beta_enabled", Op: flags.OpExists}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"beta_enabled": true},
	})
	require.NotNil(t, result)

	result = matcher(rules, flags.EvalContext{
		Attrs: map[string]any{},
	})
	assert.Nil(t, result)
}

func TestRuleMatcher_OpStartsWith(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "starts-with-rule",
			Conditions: []flags.Condition{{Attr: "email", Op: flags.OpStartsWith, Value: "@internal."}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"email": "@internal.company.com"},
	})
	require.NotNil(t, result)

	result = matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"email": "user@external.com"},
	})
	assert.Nil(t, result)
}

func TestRuleMatcher_MultipleConditions_AllMustMatch(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID: "multi-condition",
			Conditions: []flags.Condition{
				{Attr: "plan", Op: flags.OpEquals, Value: "premium"},
				{Attr: "country", Op: flags.OpEquals, Value: "US"},
				{Attr: "verified", Op: flags.OpExists},
			},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"plan": "premium", "country": "US", "verified": true},
	})
	require.NotNil(t, result)

	result = matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"plan": "premium", "country": "US"},
	})
	assert.Nil(t, result)

	result = matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"plan": "free", "country": "US", "verified": true},
	})
	assert.Nil(t, result)
}

func TestRuleMatcher_BuiltInAttrs(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "user-rule",
			Conditions: []flags.Condition{{Attr: "user_id", Op: flags.OpEquals, Value: "user-123"}},
		},
		{
			ID:         "tenant-rule",
			Conditions: []flags.Condition{{Attr: "tenant_id", Op: flags.OpEquals, Value: "tenant-456"}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		UserID: "user-123",
	})
	require.NotNil(t, result)
	assert.Equal(t, "user-rule", result.ID)

	result = matcher(rules, flags.EvalContext{
		TenantID: "tenant-456",
	})
	require.NotNil(t, result)
	assert.Equal(t, "tenant-rule", result.ID)
}

func TestRuleMatcher_UnknownOp(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "unknown-op",
			Conditions: []flags.Condition{{Attr: "x", Op: "unknown_op", Value: "y"}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"x": "y"},
	})
	assert.Nil(t, result)
}

func TestRuleMatcher_NilAttrs(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "custom-attr-rule",
			Conditions: []flags.Condition{{Attr: "custom", Op: flags.OpEquals, Value: "value"}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: nil,
	})
	assert.Nil(t, result)
}

func TestRuleMatcher_OpIn_InvalidListType(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "invalid-in",
			Conditions: []flags.Condition{{Attr: "plan", Op: flags.OpIn, Value: "not-a-slice"}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"plan": "premium"},
	})
	assert.Nil(t, result)
}

func TestRuleMatcher_OpNotIn_InvalidListType(t *testing.T) {
	t.Parallel()

	matcher := flags.DefaultRuleMatcher()

	rules := []flags.Rule{
		{
			ID:         "invalid-not-in",
			Conditions: []flags.Condition{{Attr: "plan", Op: flags.OpNotIn, Value: "not-a-slice"}},
		},
	}

	result := matcher(rules, flags.EvalContext{
		Attrs: map[string]any{"plan": "premium"},
	})
	require.NotNil(t, result)
	assert.Equal(t, "invalid-not-in", result.ID)
}
