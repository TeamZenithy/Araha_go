package volume

import (
	"strconv"

	"github.com/TeamZenithy/Araha/extensions/embed"
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
			Description:          &handler.Description{ReqPermsission: "SPEAK", Usage: "volume [0 ~ 200]"},
		},
	)
}

const (
	commandName = "volume"
	commandArg  = "vol"
)

func run(ctx handler.CommandContext) error {
	e := embed.New(ctx.Session, ctx.Message.ChannelID)
	if ms, ok := model.Music[ctx.Message.GuildID]; ok {
		guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
		if err != nil {
			return nil
		}

		if isInVoice := utils.IsInVoiceWithMusic(guild, ctx.Message.Author.ID); !isInVoice {
			e.SendEmbed(embed.BADREQ, ctx.T("music:BRNotPlaying"))
			return nil
		}
		if ctx.Arguments[commandArg] == "" {
			e.SendEmbed(embed.INFO, ctx.T("music:CurVolume", strconv.Itoa(ms.Player.GetVolume())))
			return nil
		}
		vol, err := strconv.Atoi(ctx.Arguments[commandArg])
		if err != nil {
			e.SendEmbed(embed.BADREQ, ctx.T("music:BRNumber"))
			logger.Warn(err.Error())
			return nil
		}
		if vol < 1 || vol > 200 {
			e.SendEmbed(embed.BADREQ, ctx.T("music:BRNumber"))
			return nil
		}
		prevVol := ms.Player.GetVolume()
		err = ms.Player.Volume(vol)
		if err != nil {
			logger.Warn(err.Error())
		}
		e.SendEmbed(embed.INFO, ctx.T("music:Volume", strconv.Itoa(prevVol), strconv.Itoa(ms.Player.GetVolume())))
	} else {
		e.SendEmbed(embed.BADREQ, ctx.T("music:BRNotPlaying"))
	}
	return nil
}
