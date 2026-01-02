package flags

import (
	"context"
	"time"
)

type Service struct {
	repo        Repository
	ruleMatcher RuleMatcher
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:        repo,
		ruleMatcher: DefaultRuleMatcher(),
	}
}

func NewServiceWithMatcher(repo Repository, matcher RuleMatcher) *Service {
	return &Service{
		repo:        repo,
		ruleMatcher: matcher,
	}
}

func (s *Service) Create(ctx context.Context, flag *Flag) error {
	flag.UpdatedAt = time.Now()

	return s.repo.Create(ctx, flag)
}

func (s *Service) Evaluate(ctx context.Context, key string, evalCtx EvalContext) (*EvalResult, error) {
	flag, err := s.repo.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	result := &EvalResult{
		FlagKey:     key,
		EvaluatedAt: time.Now(),
	}

	if !flag.Enabled {
		result.Value = flag.DefaultValue
		result.Reason = ReasonDisabled

		return result, nil
	}

	if rule := s.ruleMatcher(flag.Rules, evalCtx); rule != nil {
		result.Value = rule.Value
		result.Reason = ReasonRuleMatch
		result.RuleID = rule.ID

		return result, nil
	}

	result.Value = flag.DefaultValue
	result.Reason = ReasonDefault

	return result, nil
}
