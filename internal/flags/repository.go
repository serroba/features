package flags

import (
	"context"
	"errors"
)

var (
	ErrFlagNotFound = errors.New("flag not found")
	ErrFlagExists   = errors.New("flag already exists")
)

type Repository interface {
	Get(ctx context.Context, key FlagKey) (Flag, error)
	Create(ctx context.Context, flag Flag) error
}
