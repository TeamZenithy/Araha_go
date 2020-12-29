package stop

import (
	"fmt"

	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/utils"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"quit"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Usage:                map[string]string{"Required Permission": "**``SPEAK``**", "Description": "``Stop the player and clear the queue``", "Usage": fmt.Sprintf("```css\n%sstop```", utils.Prefix)},
		},
	)
}

const (
	commandName = "stop"
	commandArg  = "none"
)

func run(ctx handler.CommandContext) error {
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil
	}

	if utils.IsInVoiceWithMusic(guild, ctx.Message.Author.ID) {
		if returnedMessage := utils.LeaveAndDestroy(ctx.Session, ctx.Message.GuildID); returnedMessage != "" {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:ErrStopped")+"\n"+returnedMessage)
		} else {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:Stopped"))
		}
	} else {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:BRNotPlaying"))
	}
	return nil
}
