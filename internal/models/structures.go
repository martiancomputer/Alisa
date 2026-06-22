package models

import "github.com/bwmarrin/discordgo"

// Case represents a standardized moderation action.
type Case struct {
	ID              int64  `json:"case_id" db:"case_id"`
	UserID          string `json:"user_id" db:"user_id"`
	ModeratorID     string `json:"moderator_id" db:"moderator_id"`
	ActionType      string `json:"action_type" db:"action_type"` // e.g., "TIMEOUT", "BAN"
	Reason          string `json:"reason" db:"reason"`
	DurationSeconds *int64 `json:"duration_seconds" db:"duration_seconds"`
	Timestamp       int64  `json:"timestamp" db:"timestamp"`
}

// Rule defines an automation sequence.
type Rule struct {
	ID          int64       `json:"rule_id" db:"rule_id"`
	Name        string      `json:"name" db:"name"`
	TriggerType string      `json:"trigger_type" db:"trigger_type"` // e.g., "MESSAGE_CREATE"
	IsActive    bool        `json:"is_active" db:"is_active"`
	Conditions  []Condition `json:"conditions"`
	Actions     []Action    `json:"actions"`
}

// Condition represents a logical check required before rule execution.
type Condition struct {
	ID             int64  `json:"condition_id" db:"condition_id"`
	RuleID         int64  `json:"rule_id" db:"rule_id"`
	ConditionType  string `json:"condition_type" db:"condition_type"` // e.g., "CHANNEL", "REGEX"
	Operator       string `json:"operator" db:"operator"`             // e.g., "EQUALS", "MATCHES"
	Value          string `json:"value" db:"value"`
}

// Action represents the outcome executed if conditions are met.
type Action struct {
	ID          int64  `json:"action_id" db:"action_id"`
	RuleID      int64  `json:"rule_id" db:"rule_id"`
	ActionType  string `json:"action_type" db:"action_type"` // e.g., "DELETE_MESSAGE"
	ActionValue string `json:"action_value" db:"action_value"`
}

// EventContext encapsulates the necessary state for the automation engine to evaluate a trigger.
type EventContext struct {
	Session *discordgo.Session
	Message *discordgo.Message
	Member  *discordgo.Member
	GuildID string
}