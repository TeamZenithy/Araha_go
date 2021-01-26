package music

import (
	"fmt"
	"strconv"

	. "github.com/TeamZenithy/Araha/config"
	"github.com/TeamZenithy/Araha/extensions/embed"
	m "github.com/TeamZenithy/Araha/middlewares"
	"github.com/TeamZenithy/Araha/model"
	"github.com/bwmarrin/discordgo"

	"github.com/TeamZenithy/Araha/handler"
)

func Queue() *handler.Cmd {
	return &handler.Cmd{
		Run:         queue,
		Name:        "queue",
		Category:    "music",
		Aliases:     []string{"q"},
		Args:        []string{},
		Middlewares: []handler.HandlerFunc{m.LoadQueue(), m.VoiceWithMusic()},
		Usage:       "",
	}
}

func queue(c *handler.Context) {
	ms := c.Get("queue").(*model.MusicStruct)
	fields := []*discordgo.MessageEmbedField{}

	queueLink := "https://araha.b1ackange1.me/" + c.Locale + "/queue/" + c.Msg.GuildID
	if !Config().Release {
		queueLink = "http://localhost:8096/" + c.Locale + "/queue/" + c.Msg.GuildID
	}
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
			name += fmt.Sprintf(" (%s)", c.T("music:NowPlaying"))
		}
		user, _ := c.Session.User(song.Requester)
		name += " - " + c.T("music:ReqBy", user.Username, user.Discriminator)
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  fmt.Sprintf("[%s](%s)", song.Track.Info.Title, song.Track.Info.URI),
			Inline: false,
		})
	}
	if appendMoreField {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   c.T("music:MoreSongName"),
			Value:  c.T("music:MoreSong", queueLink),
			Inline: false,
		})
	}
	c.Embed.SendEmbed(embed.BADREQ, "", embed.AddTitle(c.T("music:MusicQueue")), embed.AddLink(queueLink), embed.AddFields(fields))
}
