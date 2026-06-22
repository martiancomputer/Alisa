package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/martiancomputer/Alisa/internal/automation"
	"github.com/martiancomputer/Alisa/internal/models"
)

// MessageHandler encapsulates the dependencies required for message event processing.
type MessageHandler struct {
	Engine *automation.Engine
}

// OnMessageCreate is bound to the discordgo session to intercept new messages.
func (h *MessageHandler) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot traffic to prevent recursion loops
	if m.Author.Bot {
		return
	}

	// Construct the context packet
	ctx := models.EventContext{
		Session: s,
		Message: m.Message,
		Member:  m.Member,
		GuildID: m.GuildID,
	}

	// Dispatch to automation engine for immediate rule processing
	h.Engine.ProcessEvent("MESSAGE_CREATE", ctx)
}