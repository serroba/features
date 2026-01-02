package flags

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, flag Flag) (Flag, error) {
	flag.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, flag); err != nil {
		return Flag{}, err
	}

	return flag, nil
}

func (s *Service) Evaluate(ctx context.Context, key FlagKey, evalCtx EvalContext) (EvalResult, error) {
	flag, err := s.repo.Get(ctx, key)
	if err != nil {
		return EvalResult{}, err
	}

	return flag.Evaluate(evalCtx), nil
}
