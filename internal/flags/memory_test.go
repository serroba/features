package flags_test

import (
	"context"
	"testing"
	"time"

	"github.com/serroba/features/internal/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryRepository_Create(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "test-flag",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(false),
		UpdatedAt:    time.Now(),
	}

	err := repo.Create(ctx, flag)
	require.NoError(t, err)
}

func TestMemoryRepository_Create_Duplicate(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	ctx := context.Background()

	flag := flags.Flag{
		Key:     "test-flag",
		Type:    flags.FlagBool,
		Enabled: true,
	}

	require.NoError(t, repo.Create(ctx, flag))
	err := repo.Create(ctx, flag)

	assert.ErrorIs(t, err, flags.ErrFlagExists)
}

func TestMemoryRepository_Get(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	ctx := context.Background()

	flag := flags.Flag{
		Key:          "test-flag",
		Type:         flags.FlagBool,
		Enabled:      true,
		DefaultValue: flags.BoolValue(true),
	}

	require.NoError(t, repo.Create(ctx, flag))

	got, err := repo.Get(ctx, "test-flag")
	require.NoError(t, err)
	assert.Equal(t, flag.Key, got.Key)
	assert.Equal(t, flag.Enabled, got.Enabled)
}

func TestMemoryRepository_Get_NotFound(t *testing.T) {
	t.Parallel()

	repo := flags.NewMemoryRepository()
	ctx := context.Background()

	_, err := repo.Get(ctx, "nonexistent")

	assert.ErrorIs(t, err, flags.ErrFlagNotFound)
}
