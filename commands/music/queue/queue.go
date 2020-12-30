package queue

import (
	"fmt"
	"strconv"

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
			Aliases:              []string{"q"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Description:          &handler.Description{ReqPermsission: "SPEAK", Usage: "queue"},
		},
	)
}

const (
	commandName = "queue"
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
		fields := []*discordgo.MessageEmbedField{}

		// ! Change This
		queueLink := "https://araha.b1ackange1.me/" + ctx.Locale + "/queue/" + ctx.Message.GuildID
		lenQueue := len(ms.Queue)
		loadLen := lenQueue
		appendMoreField := false
		if lenQueue > 5 {
			loadLen = 5
			appendMoreField = true
		}
		for i, song := range ms.Queue[:loadLen] {
			name := strconv.Itoa(i+1) + "."
			if i == 0 {
				name += fmt.Sprintf(" (%s)", ctx.T("music:NowPlaying"))
			}
			user, _ := ctx.Session.User(song.Requester)
			name += " - " + ctx.T("music:ReqBy", user.Username, user.Discriminator)
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   name,
				Value:  fmt.Sprintf("[%s](%s)", song.Track.Info.Title, song.Track.Info.URI),
				Inline: false,
			})
		}
		if appendMoreField {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   ctx.T("music:MoreSongName"),
				Value:  ctx.T("music:MoreSong", queueLink),
				Inline: false,
			})
		}
		e.SendEmbed(embed.BADREQ, "", embed.AddTitle(ctx.T("music:MusicQueue")), embed.AddLink(queueLink), embed.AddFields(fields))
	} else {
		e.SendEmbed(embed.BADREQ, ctx.T("music:NoMusic"))
	}
	return nil
}
