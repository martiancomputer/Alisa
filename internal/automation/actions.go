package automation

import (
	"log"
	"strconv"
	"time"

	"github.com/martiancomputer/Alisa/internal/models"
)

// executeDeleteMessage removes the offending message from the channel.
func executeDeleteMessage(ctx models.EventContext) {
	if ctx.Message == nil || ctx.Session == nil {
		return
	}
	err := ctx.Session.ChannelMessageDelete(ctx.Message.ChannelID, ctx.Message.ID)
	if err != nil {
		log.Printf("ERR: Action DELETE_MESSAGE failed: %v", err)
	}
}

// executeTimeout applies a communication restriction to the user.
func executeTimeout(ctx models.EventContext, durationStr string) {
	if ctx.Message == nil || ctx.Session == nil {
		return
	}

	durationSecs, err := strconv.ParseInt(durationStr, 10, 64)
	if err != nil {
		log.Printf("ERR: Action TIMEOUT failed to parse duration: %v", err)
		return
	}

	until := time.Now().Add(time.Duration(durationSecs) * time.Second)
	err = ctx.Session.GuildMemberTimeout(ctx.GuildID, ctx.Message.Author.ID, &until)
	if err != nil {
		log.Printf("ERR: Action TIMEOUT API call failed: %v", err)
	}
}