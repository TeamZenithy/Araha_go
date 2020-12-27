package nowplaying

import (
	"fmt"
	"strings"
	"time"

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
			Usage:                map[string]string{"Required Permission": "**``SPEAK``**", "Description": "``Shows current playing song's information``", "Usage": fmt.Sprintf("```css\n%snowplaying```", utils.Prefix)},
		},
	)
}

const (
	commandName = "nowplaying"
	commandArg  = "none"
)

func run(ctx handler.CommandContext) error {
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
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "You are not in a voice channel.")
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
		embed := &discordgo.MessageEmbed{
			Title:       "Now Playing",
			Description: fmt.Sprintf("[%s](%s)\nUploaded by %s", ms.Queue[0].Track.Info.Title, ms.Queue[0].Track.Info.URI, ms.Queue[0].Track.Info.Author),
			Fields:      []*discordgo.MessageEmbedField{},
			Timestamp:   time.Now().Format(time.RFC3339),
			Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: imgURL, Width: 1280, Height: 720},
			Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by %s#%s", user.Username, user.Discriminator), IconURL: user.AvatarURL("")},
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Duration",
			Value:  strings.Replace(strings.Replace(strings.Replace(time.Until(time.Now().Add(time.Duration(ms.Queue[0].Track.Info.Length)*time.Millisecond)).Round(time.Second).String(), "h", "h ", -1), "m", "m ", -1), "s", "s ", -1),
			Inline: true,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "ETA",
			Value:  strings.Replace(strings.Replace(strings.Replace(time.Until(time.Now().Add(time.Duration(ms.Queue[0].Track.Info.Length-ms.Player.Position())*time.Millisecond)).Round(time.Second).String(), "h", "h ", -1), "m", "m ", -1), "s", "s ", -1),
			Inline: true,
		})
		_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	} else {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "There is no music playing in this server.")
	}
	return nil
}
