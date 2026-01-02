package flags_test

import (
	"context"
	"testing"

	"github.com/serroba/features/internal/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := &flags.Flag{
		Key:          "new-feature",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(false),
	}

	err := svc.Create(ctx, flag)
	require.NoError(t, err)

	assert.Equal(t, int64(1), flag.Version)
	assert.False(t, flag.UpdatedAt.IsZero())
}

func TestService_Create_Duplicate(t *testing.T) {
	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := &flags.Flag{
		Key:     "new-feature",
		Type:    flags.FlagBool,
		Enabled: true,
	}

	_ = svc.Create(ctx, flag)
	err := svc.Create(ctx, flag)

	assert.ErrorIs(t, err, flags.ErrFlagExists)
}

func TestService_Evaluate_ReturnsDefault(t *testing.T) {
	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := &flags.Flag{
		Key:          "my-flag",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(true),
	}
	_ = svc.Create(ctx, flag)

	result, err := svc.Evaluate(ctx, "my-flag", flags.EvalContext{})
	require.NoError(t, err)

	assert.Equal(t, "my-flag", result.FlagKey)
	assert.Equal(t, flags.ReasonDefault, result.Reason)
	assert.Equal(t, flags.FlagBool, result.Value.Kind)
	assert.True(t, *result.Value.Bool)
}

func TestService_Evaluate_DisabledFlag(t *testing.T) {
	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	flag := &flags.Flag{
		Key:          "disabled-flag",
		Type:         flags.FlagBool,
		Enabled:      false,
		DefaultValue: flags.BoolValue(false),
	}
	_ = svc.Create(ctx, flag)

	result, err := svc.Evaluate(ctx, "disabled-flag", flags.EvalContext{})
	require.NoError(t, err)

	assert.Equal(t, flags.ReasonDisabled, result.Reason)
	assert.False(t, *result.Value.Bool)
}

func TestService_Evaluate_NotFound(t *testing.T) {
	repo := flags.NewMemoryRepository()
	svc := flags.NewService(repo)
	ctx := context.Background()

	_, err := svc.Evaluate(ctx, "nonexistent", flags.EvalContext{})

	assert.ErrorIs(t, err, flags.ErrFlagNotFound)
}
