package nowplaying

import (
	"strings"
	"time"

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
			Aliases:              []string{"np"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Description:          &handler.Description{ReqPermsission: "SPEAK", Usage: "nowplaying"},
		},
	)
}

const (
	commandName = "nowplaying"
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
		user, _ := ctx.Session.User(ms.Queue[0].Requester)
		var imgURL = ""
		if strings.Contains(ms.Queue[0].Track.Info.URI, "youtube.com") || strings.Contains(ms.Queue[0].Track.Info.URI, "yt.be") {
			imgURL, _ = utils.GetYTThumbnail(ms.Queue[0].Track.Info.URI)
		}
		if strings.Contains(ms.Queue[0].Track.Info.URI, "soundcloud.com") {
			imgURL, _ = utils.GetSCThumbnail(ms.Queue[0].Track.Info.URI)
		}
		description := ctx.T("music:SongInfo", ms.Queue[0].Track.Info.Title, ms.Queue[0].Track.Info.URI, ms.Queue[0].Track.Info.Author)
		fields := []*discordgo.MessageEmbedField{
			{
				Name:   ctx.T("music:Duration"),
				Value:  utils.TimeFormat(time.Duration(ms.Queue[0].Track.Info.Length)),
				Inline: true,
			},
			{
				Name:   ctx.T("music:ETA"),
				Value:  utils.TimeFormat(time.Duration(ms.Queue[0].Track.Info.Length - ms.Player.Position())),
				Inline: true,
			},
		}
		e.SendEmbed(embed.INFO, description, embed.AddTitle("Now Playing"), embed.AddThumbnail(imgURL), embed.AddProfileFooter(user, ctx.T), embed.AddFields(fields))
	} else {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:NoMusic"))
	}
	return nil
}
