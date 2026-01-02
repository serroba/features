package handler

import (
	"github.com/serroba/features/internal/flags"
)

func ToFlag(body CreateFlagBody) flags.Flag {
	return flags.Flag{
		Key:          flags.FlagKey(body.Key),
		Type:         flags.FlagType(body.Type),
		Enabled:      body.Enabled,
		DefaultValue: toValue(body.DefaultValue),
		Rules:        toRules(body.Rules),
	}
}

func toRules(bodies []RuleBody) []flags.Rule {
	if len(bodies) == 0 {
		return nil
	}

	rules := make([]flags.Rule, len(bodies))
	for i, b := range bodies {
		rules[i] = flags.Rule{
			ID:         b.ID,
			Conditions: toConditions(b.Conditions),
			Value:      toValue(b.Value),
		}
	}

	return rules
}

func toConditions(bodies []ConditionBody) []flags.Condition {
	if len(bodies) == 0 {
		return nil
	}

	conditions := make([]flags.Condition, len(bodies))
	for i, b := range bodies {
		conditions[i] = flags.Condition{
			Attr:  b.Attr,
			Op:    flags.ConditionOp(b.Op),
			Value: b.Value,
		}
	}

	return conditions
}

func toValue(body ValueBody) flags.Value {
	return flags.Value{
		Kind:   flags.FlagType(body.Kind),
		Bool:   body.Bool,
		String: body.String,
		Number: body.Number,
	}
}

func ToEvalContext(body EvaluateFlagBody) flags.EvalContext {
	return flags.EvalContext{
		TenantID: body.TenantID,
		UserID:   body.UserID,
		Attrs:    body.Attrs,
	}
}

func ToEvalResultBody(result flags.EvalResult) EvalResultBody {
	return EvalResultBody{
		FlagKey:     result.FlagKey,
		Value:       toValueBody(result.Value),
		Reason:      string(result.Reason),
		RuleID:      result.RuleID,
		EvaluatedAt: result.EvaluatedAt,
	}
}

func toValueBody(value flags.Value) ValueBody {
	return ValueBody{
		Kind:   string(value.Kind),
		Bool:   value.Bool,
		String: value.String,
		Number: value.Number,
	}
}
