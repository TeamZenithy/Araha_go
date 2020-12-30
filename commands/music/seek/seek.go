package seek

import (
	"strconv"

	"github.com/TeamZenithy/Araha/extensions/embed"
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
			Description:          &handler.Description{ReqPermsission: "SPEAK", Usage: "seek [second]"},
		},
	)
}

const (
	commandName = "seek"
	commandArg  = "second"
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
		pos, errNotSecond := strconv.Atoi(ctx.Arguments[commandArg])
		pos = pos * 1000
		if errNotSecond != nil || pos < 0 {
			e.SendEmbed(embed.BADREQ, ctx.T("music:BRTime"))
			return nil
		}
		if ms.Queue[0].Track.Info.Length <= int64(pos) {
			e.SendEmbed(embed.BADREQ, ctx.T("music:BRTime2"))
			return nil
		}
		ms.Player.Seek(int64(pos))
		e.SendEmbed(embed.BADREQ, ctx.T("music:SeekTo", ctx.Arguments[commandArg]))
	} else {
		e.SendEmbed(embed.BADREQ, ctx.T("music:NoMusic"))
	}
	return nil
}
