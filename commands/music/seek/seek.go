package seek

import (
	"fmt"
	"strconv"

	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Usage:                map[string]string{"Required Permission": "**``SPEAK``**", "Description": "``Navigate to the requested location of the song that is currently playing.``", "Usage": fmt.Sprintf("```css\n%sseek [second]```", utils.Prefix)},
		},
	)
}

const (
	commandName = "seek"
	commandArg  = "second"
)

func run(ctx handler.CommandContext) error {
	if ms, ok := model.Music[ctx.Message.GuildID]; ok {
		guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
		if err != nil {
			return nil
		}

		if isInVoice := utils.IsInVoiceWithMusic(guild, ctx.Message.Author.ID); !isInVoice {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "You're not listening to my music :(")
			return nil
		}
		pos, errNotSecond := strconv.Atoi(ctx.Arguments[commandArg])
		pos = pos * 1000
		if errNotSecond != nil || pos < 0 {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Please give me valid time")
			return nil
		}
		if ms.Queue[0].Track.Info.Length <= int64(pos) {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Requested time is same with song or longer")
			return nil
		}
		ms.Player.Seek(int64(pos))
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(":stopwatch: Moved to %s.", ctx.Arguments[commandArg]))
	} else {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "There is no music playing.")
	}
	return nil
}
