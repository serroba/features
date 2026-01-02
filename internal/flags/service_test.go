package flags_test

import (
	"context"
	"testing"

	"github.com/serroba/features/internal/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	input := flags.Flag{
		Key:          "new-feature",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(false),
	}

	created, err := svc.Create(ctx, input)
	require.NoError(t, err)

	assert.Equal(t, flags.FlagKey("new-feature"), created.Key)
	assert.False(t, created.UpdatedAt.IsZero())
}

func TestService_Create_Duplicate(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	input := flags.Flag{
		Key:     "new-feature",
		Type:    flags.FlagBool,
		Enabled: true,
	}

	_, err := svc.Create(ctx, input)
	require.NoError(t, err)

	_, err = svc.Create(ctx, input)
	assert.ErrorIs(t, err, flags.ErrFlagExists)
}

func TestService_Evaluate_ReturnsDefault(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "my-flag",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(true),
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "my-flag", flags.EvalContext{})
	require.NoError(t, err)

	assert.Equal(t, flags.FlagKey("my-flag"), result.FlagKey)
	assert.Equal(t, flags.ReasonDefault, result.Reason)
	assert.Equal(t, flags.FlagBool, result.Value.Kind)
	assert.True(t, *result.Value.Bool)
}

func TestService_Evaluate_DisabledFlag(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "disabled-flag",
		Type:         flags.FlagBool,
		Enabled:      false,
		DefaultValue: flags.BoolValue(false),
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "disabled-flag", flags.EvalContext{})
	require.NoError(t, err)

	assert.Equal(t, flags.ReasonDisabled, result.Reason)
	assert.False(t, *result.Value.Bool)
}

func TestService_Evaluate_NotFound(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	_, err := svc.Evaluate(ctx, "nonexistent", flags.EvalContext{})

	assert.ErrorIs(t, err, flags.ErrFlagNotFound)
}

func TestService_Evaluate_RuleMatch_Equals(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "premium-feature",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(false),
		Rules: []flags.Rule{
			{
				ID: "premium-users",
				Conditions: []flags.Condition{
					{Attr: "plan", Op: flags.OpEquals, Value: "premium"},
				},
				Value: flags.BoolValue(true),
			},
		},
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "premium-feature", flags.EvalContext{
		Attrs: map[string]any{"plan": "premium"},
	})
	require.NoError(t, err)

	assert.Equal(t, flags.ReasonRuleMatch, result.Reason)
	assert.Equal(t, "premium-users", result.RuleID)
	assert.True(t, *result.Value.Bool)
}

func TestService_Evaluate_RuleMatch_NoMatch(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "premium-feature",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(false),
		Rules: []flags.Rule{
			{
				ID: "premium-users",
				Conditions: []flags.Condition{
					{Attr: "plan", Op: flags.OpEquals, Value: "premium"},
				},
				Value: flags.BoolValue(true),
			},
		},
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "premium-feature", flags.EvalContext{
		Attrs: map[string]any{"plan": "free"},
	})
	require.NoError(t, err)

	assert.Equal(t, flags.ReasonDefault, result.Reason)
	assert.False(t, *result.Value.Bool)
}

func TestService_Evaluate_RuleMatch_In(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "beta-feature",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(false),
		Rules: []flags.Rule{
			{
				ID: "beta-tenants",
				Conditions: []flags.Condition{
					{Attr: "tenant_id", Op: flags.OpIn, Value: []any{"tenant-1", "tenant-2"}},
				},
				Value: flags.BoolValue(true),
			},
		},
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "beta-feature", flags.EvalContext{
		TenantID: "tenant-2",
	})
	require.NoError(t, err)

	assert.Equal(t, flags.ReasonRuleMatch, result.Reason)
	assert.True(t, *result.Value.Bool)
}

func TestService_Evaluate_RuleMatch_MultipleConditions(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "geo-feature",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(false),
		Rules: []flags.Rule{
			{
				ID: "premium-us",
				Conditions: []flags.Condition{
					{Attr: "plan", Op: flags.OpEquals, Value: "premium"},
					{Attr: "country", Op: flags.OpEquals, Value: "US"},
				},
				Value: flags.BoolValue(true),
			},
		},
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "geo-feature", flags.EvalContext{
		Attrs: map[string]any{"plan": "premium", "country": "US"},
	})
	require.NoError(t, err)
	assert.Equal(t, flags.ReasonRuleMatch, result.Reason)

	result, err = svc.Evaluate(ctx, "geo-feature", flags.EvalContext{
		Attrs: map[string]any{"plan": "premium", "country": "UK"},
	})
	require.NoError(t, err)
	assert.Equal(t, flags.ReasonDefault, result.Reason)
}

func TestService_Evaluate_RuleMatch_FirstMatchWins(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "tiered-feature",
		Type:         flags.FlagString,
		Enabled:      true,
		DefaultValue: flags.StringValue("basic"),
		Rules: []flags.Rule{
			{
				ID: "enterprise",
				Conditions: []flags.Condition{
					{Attr: "plan", Op: flags.OpEquals, Value: "enterprise"},
				},
				Value: flags.StringValue("full"),
			},
			{
				ID: "premium",
				Conditions: []flags.Condition{
					{Attr: "plan", Op: flags.OpEquals, Value: "premium"},
				},
				Value: flags.StringValue("partial"),
			},
		},
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "tiered-feature", flags.EvalContext{
		Attrs: map[string]any{"plan": "enterprise"},
	})
	require.NoError(t, err)

	assert.Equal(t, "enterprise", result.RuleID)
	assert.Equal(t, "full", *result.Value.String)
}

func TestService_Evaluate_RuleMatch_StartsWith(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "internal-feature",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(false),
		Rules: []flags.Rule{
			{
				ID: "internal-emails",
				Conditions: []flags.Condition{
					{Attr: "email", Op: flags.OpStartsWith, Value: "admin@"},
				},
				Value: flags.BoolValue(true),
			},
		},
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "internal-feature", flags.EvalContext{
		Attrs: map[string]any{"email": "admin@company.com"},
	})
	require.NoError(t, err)
	assert.Equal(t, flags.ReasonRuleMatch, result.Reason)

	result, err = svc.Evaluate(ctx, "internal-feature", flags.EvalContext{
		Attrs: map[string]any{"email": "user@company.com"},
	})
	require.NoError(t, err)
	assert.Equal(t, flags.ReasonDefault, result.Reason)
}

func TestService_Evaluate_StringValue(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "welcome-message",
		Type:         flags.FlagString,
		Enabled:      true,
		DefaultValue: flags.StringValue("Hello, user!"),
		Rules: []flags.Rule{
			{
				ID: "vip-message",
				Conditions: []flags.Condition{
					{Attr: "tier", Op: flags.OpEquals, Value: "vip"},
				},
				Value: flags.StringValue("Welcome back, VIP!"),
			},
		},
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "welcome-message", flags.EvalContext{
		Attrs: map[string]any{"tier": "vip"},
	})
	require.NoError(t, err)
	assert.Equal(t, flags.ReasonRuleMatch, result.Reason)
	assert.Equal(t, flags.FlagString, result.Value.Kind)
	assert.Equal(t, "Welcome back, VIP!", *result.Value.String)

	result, err = svc.Evaluate(ctx, "welcome-message", flags.EvalContext{
		Attrs: map[string]any{"tier": "regular"},
	})
	require.NoError(t, err)
	assert.Equal(t, flags.ReasonDefault, result.Reason)
	assert.Equal(t, "Hello, user!", *result.Value.String)
}

func TestService_WithCustomMatcher(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	customMatcher := func(_ []flags.Rule, _ flags.EvalContext) (flags.Rule, bool) {
		return flags.Rule{
			ID:    "custom-rule",
			Value: flags.StringValue("custom-value"),
		}, true
	}
	svc := flags.NewServiceWithMatcher(repo, customMatcher)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "test-flag",
		Type:         flags.FlagString,
		Enabled:      true,
		DefaultValue: flags.StringValue("default"),
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "test-flag", flags.EvalContext{})
	require.NoError(t, err)

	assert.Equal(t, flags.ReasonRuleMatch, result.Reason)
	assert.Equal(t, "custom-rule", result.RuleID)
	assert.Equal(t, "custom-value", *result.Value.String)
}

func TestService_Evaluate_NumberValue(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "rate-limit",
		Type:         flags.FlagNumber,
		Enabled:      true,
		DefaultValue: flags.NumberValue(100),
		Rules: []flags.Rule{
			{
				ID: "premium-limit",
				Conditions: []flags.Condition{
					{Attr: "plan", Op: flags.OpEquals, Value: "premium"},
				},
				Value: flags.NumberValue(1000),
			},
			{
				ID: "enterprise-limit",
				Conditions: []flags.Condition{
					{Attr: "plan", Op: flags.OpEquals, Value: "enterprise"},
				},
				Value: flags.NumberValue(10000),
			},
		},
	}
	_, err := svc.Create(ctx, flag)
	require.NoError(t, err)

	result, err := svc.Evaluate(ctx, "rate-limit", flags.EvalContext{
		Attrs: map[string]any{"plan": "enterprise"},
	})
	require.NoError(t, err)
	assert.Equal(t, flags.ReasonRuleMatch, result.Reason)
	assert.Equal(t, flags.FlagNumber, result.Value.Kind)
	assert.InDelta(t, float64(10000), *result.Value.Number, 0.001)

	result, err = svc.Evaluate(ctx, "rate-limit", flags.EvalContext{
		Attrs: map[string]any{"plan": "free"},
	})
	require.NoError(t, err)
	assert.Equal(t, flags.ReasonDefault, result.Reason)
	assert.InDelta(t, float64(100), *result.Value.Number, 0.001)
}
