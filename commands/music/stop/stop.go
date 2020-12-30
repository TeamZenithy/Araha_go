package stop

import (
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
			Aliases:              []string{"quit"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Description:          &handler.Description{ReqPermsission: "SPEAK", Usage: "stop"},
		},
	)
}

const (
	commandName = "stop"
	commandArg  = "none"
)

func run(ctx handler.CommandContext) error {
	e := embed.New(ctx.Session, ctx.Message.ChannelID)

	if _, ok := model.Music[ctx.Message.GuildID]; ok {
		guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
		if err != nil {
			return nil
		}

		if utils.IsInVoiceWithMusic(guild, ctx.Message.Author.ID) {
			if returnedMessage := utils.LeaveAndDestroy(ctx.Session, ctx.Message.GuildID); returnedMessage != "" {
				ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:ErrStopped")+"\n"+returnedMessage)
			} else {
				e.SendEmbed(embed.INFO, ctx.T("music:Stopped"))
				return nil
			}
		} else {
			e.SendEmbed(embed.BADREQ, ctx.T("music:BRNotPlaying"))
			return nil
		}
	} else {
		e.SendEmbed(embed.BADREQ, ctx.T("music:NoMusic"))
	}

	return nil
}
