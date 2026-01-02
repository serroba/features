package handler_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/serroba/features/internal/flags"
	"github.com/serroba/features/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandler_CreateFlag(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockService := NewMockFlagService(ctrl)
	h := handler.New(mockService)
	ctx := context.Background()

	mockService.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, flag flags.Flag) (flags.Flag, error) {
			flag.UpdatedAt = time.Now()

			return flag, nil
		})

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
	assert.False(t, resp.Body.CreatedAt.IsZero())
}

func TestHandler_CreateFlag_Duplicate(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockService := NewMockFlagService(ctrl)
	h := handler.New(mockService)
	ctx := context.Background()

	mockService.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(flags.Flag{}, flags.ErrFlagExists)

	req := &handler.CreateFlagRequest{
		Body: handler.CreateFlagBody{
			Key:     "test-flag",
			Type:    "bool",
			Enabled: true,
		},
	}

	_, err := h.CreateFlag(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestHandler_CreateFlag_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockService := NewMockFlagService(ctrl)
	h := handler.New(mockService)
	ctx := context.Background()

	mockService.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(flags.Flag{}, errors.New("database connection failed"))

	req := &handler.CreateFlagRequest{
		Body: handler.CreateFlagBody{
			Key:     "test-flag",
			Type:    "bool",
			Enabled: true,
		},
	}

	_, err := h.CreateFlag(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create flag")
}

func TestHandler_EvaluateFlag(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockService := NewMockFlagService(ctrl)
	h := handler.New(mockService)
	ctx := context.Background()

	boolVal := true
	mockService.EXPECT().
		Evaluate(gomock.Any(), "my-flag", gomock.Any()).
		Return(flags.EvalResult{
			FlagKey:     "my-flag",
			Value:       flags.Value{Kind: flags.FlagBool, Bool: &boolVal},
			Reason:      flags.ReasonDefault,
			EvaluatedAt: time.Now(),
		}, nil)

	req := &handler.EvaluateFlagRequest{
		Key: "my-flag",
		Body: handler.EvaluateFlagBody{
			TenantID: "tenant-1",
			UserID:   "user-1",
		},
	}

	resp, err := h.EvaluateFlag(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, "my-flag", resp.Body.FlagKey)
	assert.Equal(t, "default", resp.Body.Reason)
	assert.Equal(t, "bool", resp.Body.Value.Kind)
	assert.True(t, *resp.Body.Value.Bool)
}

func TestHandler_EvaluateFlag_NotFound(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockService := NewMockFlagService(ctrl)
	h := handler.New(mockService)
	ctx := context.Background()

	mockService.EXPECT().
		Evaluate(gomock.Any(), "nonexistent", gomock.Any()).
		Return(flags.EvalResult{}, flags.ErrFlagNotFound)

	req := &handler.EvaluateFlagRequest{
		Key:  "nonexistent",
		Body: handler.EvaluateFlagBody{},
	}

	_, err := h.EvaluateFlag(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestHandler_EvaluateFlag_RuleMatch(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockService := NewMockFlagService(ctrl)
	h := handler.New(mockService)
	ctx := context.Background()

	boolVal := true
	mockService.EXPECT().
		Evaluate(gomock.Any(), "premium-feature", gomock.Any()).
		Return(flags.EvalResult{
			FlagKey:     "premium-feature",
			Value:       flags.Value{Kind: flags.FlagBool, Bool: &boolVal},
			Reason:      flags.ReasonRuleMatch,
			RuleID:      "premium-users",
			EvaluatedAt: time.Now(),
		}, nil)

	req := &handler.EvaluateFlagRequest{
		Key: "premium-feature",
		Body: handler.EvaluateFlagBody{
			Attrs: map[string]any{"plan": "premium"},
		},
	}

	resp, err := h.EvaluateFlag(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, "rule_match", resp.Body.Reason)
	assert.Equal(t, "premium-users", resp.Body.RuleID)
	assert.True(t, *resp.Body.Value.Bool)
}

func TestHandler_EvaluateFlag_Disabled(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockService := NewMockFlagService(ctrl)
	h := handler.New(mockService)
	ctx := context.Background()

	boolVal := false
	mockService.EXPECT().
		Evaluate(gomock.Any(), "disabled-flag", gomock.Any()).
		Return(flags.EvalResult{
			FlagKey:     "disabled-flag",
			Value:       flags.Value{Kind: flags.FlagBool, Bool: &boolVal},
			Reason:      flags.ReasonDisabled,
			EvaluatedAt: time.Now(),
		}, nil)

	req := &handler.EvaluateFlagRequest{
		Key:  "disabled-flag",
		Body: handler.EvaluateFlagBody{},
	}

	resp, err := h.EvaluateFlag(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, "disabled", resp.Body.Reason)
}

func TestHandler_EvaluateFlag_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockService := NewMockFlagService(ctrl)
	h := handler.New(mockService)
	ctx := context.Background()

	mockService.EXPECT().
		Evaluate(gomock.Any(), "test-flag", gomock.Any()).
		Return(flags.EvalResult{}, errors.New("database connection failed"))

	req := &handler.EvaluateFlagRequest{
		Key:  "test-flag",
		Body: handler.EvaluateFlagBody{},
	}

	_, err := h.EvaluateFlag(ctx, req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to evaluate flag")
}
