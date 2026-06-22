package automation

import (
	"log"

	"github.com/martiancomputer/Alisa/internal/models"
)

// Engine manages the caching and execution of automation rules.
type Engine struct {
	Rules []models.Rule
	// TODO: Add database connection pointer to periodically refresh Rules cache
}

// ProcessEvent routes a standard Discord event through the active rule matrix.
func (e *Engine) ProcessEvent(triggerType string, ctx models.EventContext) {
	for _, rule := range e.Rules {
		if !rule.IsActive || rule.TriggerType != triggerType {
			continue
		}

		if e.evaluateConditions(rule.Conditions, ctx) {
			e.executeActions(rule.Actions, ctx)
		}
	}
}

// evaluateConditions processes logical constraints. Returns false if any constraint fails.
func (e *Engine) evaluateConditions(conditions []models.Condition, ctx models.EventContext) bool {
	if len(conditions) == 0 {
		return false // Failsafe: Rules without conditions should not execute blindly.
	}

	for _, cond := range conditions {
		matched := false
		switch cond.ConditionType {
		case "CHANNEL":
			if ctx.Message != nil {
				matched = checkStringMatch(ctx.Message.ChannelID, cond.Operator, cond.Value)
			}
		case "REGEX":
			if ctx.Message != nil {
				matched = checkRegexMatch(ctx.Message.Content, cond.Value)
			}
		// Additional condition logic (Account Age, Role, etc.)
		default:
			log.Printf("WARN: Unrecognized condition type encountered: %s", cond.ConditionType)
		}

		if !matched {
			return false // AND logic enforcement: all conditions must pass.
		}
	}
	return true
}

func (e *Engine) executeActions(actions []models.Action, ctx models.EventContext) {
	for _, action := range actions {
		switch action.ActionType {
		case "DELETE_MESSAGE":
			executeDeleteMessage(ctx)
		case "TIMEOUT":
			executeTimeout(ctx, action.ActionValue)
		// Additional action logic (Warn, Ban, Log Event)
		default:
			log.Printf("WARN: Unrecognized action type encountered: %s", action.ActionType)
		}
	}
}