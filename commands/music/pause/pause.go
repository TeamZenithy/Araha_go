package pause

import (
	"github.com/TeamZenithy/Araha/extensions/embed"
	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{""},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Description:          &handler.Description{ReqPermsission: "SPEAK", Usage: "pause"},
		},
	)
}

const (
	commandName = "pause"
	commandArg  = "none"
)

func run(ctx handler.CommandContext) error {
	e := embed.New(ctx.Session, ctx.Message.ChannelID)
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil
	}

	var userVoiceState discordgo.VoiceState
	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			userVoiceState = *vs
		}
	}
	if userVoiceState.UserID == "" {
		e.SendEmbed(embed.BADREQ, ctx.T("music:NotInVoiceChannel"))
		return nil
	}
	ms, ok := model.Music[ctx.Message.GuildID]
	if ok && len(ms.Queue) > 0 {
		if ms.Player.Paused() {
			e.SendEmbed(embed.BADREQ, ctx.T("music:AlreadyPaused"))
		} else {
			ms.Player.Pause(true)
			e.SendEmbed(embed.BADREQ, ctx.T("music:Paused"))
		}
	} else {
		e.SendEmbed(embed.BADREQ, ctx.T("music:NoMusic"))
	}
	return nil
}
