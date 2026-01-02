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
	return NewServiceWithMatcher(repo, DefaultRuleMatcher())
}

func NewServiceWithMatcher(repo Repository, matcher RuleMatcher) *Service {
	return &Service{
		repo:        repo,
		ruleMatcher: matcher,
	}
}

func (s *Service) Create(ctx context.Context, flag Flag) (Flag, error) {
	flag.UpdatedAt = time.Now()

	if err := s.repo.Create(ctx, flag); err != nil {
		return Flag{}, err
	}

	return flag, nil
}

func (s *Service) Evaluate(ctx context.Context, key string, evalCtx EvalContext) (EvalResult, error) {
	flag, err := s.repo.Get(ctx, key)
	if err != nil {
		return EvalResult{}, err
	}

	result := EvalResult{
		FlagKey:     key,
		EvaluatedAt: time.Now(),
	}

	if !flag.Enabled {
		result.Value = flag.DefaultValue
		result.Reason = ReasonDisabled

		return result, nil
	}

	if rule, ok := s.ruleMatcher(flag.Rules, evalCtx); ok {
		result.Value = rule.Value
		result.Reason = ReasonRuleMatch
		result.RuleID = rule.ID

		return result, nil
	}

	result.Value = flag.DefaultValue
	result.Reason = ReasonDefault

	return result, nil
}
