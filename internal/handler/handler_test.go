package handler_test

import (
	"context"
	"testing"

	"github.com/serroba/features/internal/flags"
	"github.com/serroba/features/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_CreateFlag(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	h := handler.New(svc)
	ctx := context.Background()

	boolVal := true
	req := &handler.CreateFlagRequest{
		Body: handler.CreateFlagBody{
			Key:     "test-flag",
			Type:    "bool",
			Enabled: true,
			DefaultValue: handler.ValueBody{
				Kind: "bool",
				Bool: &boolVal,
			},
		},
	}

	resp, err := h.CreateFlag(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, "test-flag", resp.Body.Key)
	assert.Equal(t, int64(1), resp.Body.Version)
	assert.False(t, resp.Body.CreatedAt.IsZero())
}

func TestHandler_CreateFlag_Duplicate(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	h := handler.New(svc)
	ctx := context.Background()

	req := &handler.CreateFlagRequest{
		Body: handler.CreateFlagBody{
			Key:     "test-flag",
			Type:    "bool",
			Enabled: true,
		},
	}

	_, err := h.CreateFlag(ctx, req)
	require.NoError(t, err)

	_, err = h.CreateFlag(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestHandler_EvaluateFlag(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	h := handler.New(svc)
	ctx := context.Background()

	boolVal := true
	createReq := &handler.CreateFlagRequest{
		Body: handler.CreateFlagBody{
			Key:     "my-flag",
			Type:    "bool",
			Enabled: true,
			DefaultValue: handler.ValueBody{
				Kind: "bool",
				Bool: &boolVal,
			},
		},
	}
	_, err := h.CreateFlag(ctx, createReq)
	require.NoError(t, err)

	evalReq := &handler.EvaluateFlagRequest{
		Key: "my-flag",
		Body: handler.EvaluateFlagBody{
			TenantID: "tenant-1",
			UserID:   "user-1",
		},
	}

	resp, err := h.EvaluateFlag(ctx, evalReq)
	require.NoError(t, err)

	assert.Equal(t, "my-flag", resp.Body.FlagKey)
	assert.Equal(t, "default", resp.Body.Reason)
	assert.Equal(t, "bool", resp.Body.Value.Kind)
	assert.True(t, *resp.Body.Value.Bool)
}

func TestHandler_EvaluateFlag_NotFound(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	h := handler.New(svc)
	ctx := context.Background()

	req := &handler.EvaluateFlagRequest{
		Key:  "nonexistent",
		Body: handler.EvaluateFlagBody{},
	}

	_, err := h.EvaluateFlag(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandler_EvaluateFlag_WithRules(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	h := handler.New(svc)
	ctx := context.Background()

	boolFalse := false
	boolTrue := true
	createReq := &handler.CreateFlagRequest{
		Body: handler.CreateFlagBody{
			Key:     "premium-feature",
			Type:    "bool",
			Enabled: true,
			DefaultValue: handler.ValueBody{
				Kind: "bool",
				Bool: &boolFalse,
			},
			Rules: []handler.RuleBody{
				{
					ID: "premium-users",
					Conditions: []handler.ConditionBody{
						{Attr: "plan", Op: "eq", Value: "premium"},
					},
					Value: handler.ValueBody{
						Kind: "bool",
						Bool: &boolTrue,
					},
				},
			},
		},
	}
	_, err := h.CreateFlag(ctx, createReq)
	require.NoError(t, err)

	evalReq := &handler.EvaluateFlagRequest{
		Key: "premium-feature",
		Body: handler.EvaluateFlagBody{
			Attrs: map[string]any{"plan": "premium"},
		},
	}

	resp, err := h.EvaluateFlag(ctx, evalReq)
	require.NoError(t, err)

	assert.Equal(t, "rule_match", resp.Body.Reason)
	assert.Equal(t, "premium-users", resp.Body.RuleID)
	assert.True(t, *resp.Body.Value.Bool)
}

func TestHandler_EvaluateFlag_Disabled(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	h := handler.New(svc)
	ctx := context.Background()

	boolVal := false
	createReq := &handler.CreateFlagRequest{
		Body: handler.CreateFlagBody{
			Key:     "disabled-flag",
			Type:    "bool",
			Enabled: false,
			DefaultValue: handler.ValueBody{
				Kind: "bool",
				Bool: &boolVal,
			},
		},
	}
	_, err := h.CreateFlag(ctx, createReq)
	require.NoError(t, err)

	evalReq := &handler.EvaluateFlagRequest{
		Key:  "disabled-flag",
		Body: handler.EvaluateFlagBody{},
	}

	resp, err := h.EvaluateFlag(ctx, evalReq)
	require.NoError(t, err)

	assert.Equal(t, "disabled", resp.Body.Reason)
}
