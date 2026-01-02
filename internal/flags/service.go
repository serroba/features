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

func (s *Service) Create(ctx context.Context, flag *Flag) error {
	flag.UpdatedAt = time.Now()

	if flag.Version == 0 {
		flag.Version = 1
	}

	return s.repo.Create(ctx, flag)
}

func (s *Service) Evaluate(ctx context.Context, key string, _ EvalContext) (*EvalResult, error) {
	flag, err := s.repo.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	result := &EvalResult{
		FlagKey:     key,
		Version:     flag.Version,
		EvaluatedAt: time.Now(),
	}

	if !flag.Enabled {
		result.Value = flag.DefaultValue
		result.Reason = ReasonDisabled

		return result, nil
	}

	result.Value = flag.DefaultValue
	result.Reason = ReasonDefault

	return result, nil
}
