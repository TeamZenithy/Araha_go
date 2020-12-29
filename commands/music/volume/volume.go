package volume

import (
	"fmt"
	"strconv"

	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"vol"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Usage:                map[string]string{"Required Permission": "**``SPEAK``**", "Description": "``Adjust volume of player.``", "Usage": fmt.Sprintf("```css\n%svolume [0 ~ 200]```", utils.Prefix)},
		},
	)
}

const (
	commandName = "volume"
	commandArg  = "vol"
)

func run(ctx handler.CommandContext) error {
	if ms, ok := model.Music[ctx.Message.GuildID]; ok {
		guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
		if err != nil {
			return nil
		}

		if isInVoice := utils.IsInVoiceWithMusic(guild, ctx.Message.Author.ID); !isInVoice {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:BRNotPlaying"))
			return nil
		}
		if ctx.Arguments[commandArg] == "" {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:CurVolume", strconv.Itoa(ms.Player.GetVolume())))
			return nil
		}
		vol, err := strconv.Atoi(ctx.Arguments[commandArg])
		if err != nil {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:BRNumber"))
			logger.Warn(err.Error())
			return nil
		}
		if vol < 1 || vol > 200 {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:BRNumber"))
			return nil
		}
		prevVol := ms.Player.GetVolume()
		err = ms.Player.Volume(vol)
		if err != nil {
			logger.Warn(err.Error())
		}
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:Volume", strconv.Itoa(prevVol), strconv.Itoa(ms.Player.GetVolume())))
	} else {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:BRNotPlaying"))
	}
	return nil
}
