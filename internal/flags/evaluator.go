package flags

import "strings"

type RuleMatcher func(rules []Rule, evalCtx EvalContext) (Rule, bool)

func DefaultRuleMatcher() RuleMatcher {
	return func(rules []Rule, evalCtx EvalContext) (Rule, bool) {
		for _, rule := range rules {
			if matchesRule(rule, evalCtx) {
				return rule, true
			}
		}

		return Rule{}, false
	}
}

func matchesRule(rule Rule, evalCtx EvalContext) bool {
	for _, cond := range rule.Conditions {
		if !matchesCondition(cond, evalCtx) {
			return false
		}
	}

	return true
}

func matchesCondition(cond Condition, evalCtx EvalContext) bool {
	attrValue := getAttrValue(cond.Attr, evalCtx)

	switch cond.Op {
	case OpEquals:
		return attrValue == cond.Value
	case OpNotEquals:
		return attrValue != cond.Value
	case OpIn:
		return valueIn(attrValue, cond.Value)
	case OpNotIn:
		return !valueIn(attrValue, cond.Value)
	case OpExists:
		return attrValue != nil
	case OpStartsWith:
		str, ok := attrValue.(string)
		prefix, prefixOk := cond.Value.(string)

		return ok && prefixOk && strings.HasPrefix(str, prefix)
	default:
		return false
	}
}

func getAttrValue(attr string, evalCtx EvalContext) any {
	switch attr {
	case "user_id":
		return evalCtx.UserID
	case "tenant_id":
		return evalCtx.TenantID
	default:
		if evalCtx.Attrs != nil {
			return evalCtx.Attrs[attr]
		}

		return nil
	}
}

func valueIn(value any, list any) bool {
	slice, ok := list.([]any)
	if !ok {
		return false
	}

	for _, item := range slice {
		if value == item {
			return true
		}
	}

	return false
}
