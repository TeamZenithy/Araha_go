package music

import (
	"strings"
	"time"

	"github.com/TeamZenithy/Araha/extensions/embed"
	m "github.com/TeamZenithy/Araha/middlewares"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"

	"github.com/TeamZenithy/Araha/handler"
)

func NowPlaying() *handler.Cmd {
	return &handler.Cmd{
		Run:         nowPlaying,
		Name:        "nowplaying",
		Category:    "music",
		Aliases:     []string{"np"},
		Args:        []string{},
		Middlewares: []handler.HandlerFunc{m.LoadQueue(), m.VoiceWithMusic()},
		Usage:       "",
	}
}

func nowPlaying(c *handler.Context) {
	ms := c.Get("queue").(*model.MusicStruct)
	user, _ := c.Session.User(ms.Queue[0].Requester)
	var imgURL = ""
	if strings.Contains(ms.Queue[0].Track.Info.URI, "youtube.com") || strings.Contains(ms.Queue[0].Track.Info.URI, "yt.be") {
		imgURL, _ = utils.GetYTThumbnail(ms.Queue[0].Track.Info.URI)
	}
	if strings.Contains(ms.Queue[0].Track.Info.URI, "soundcloud.com") {
		imgURL, _ = utils.GetSCThumbnail(ms.Queue[0].Track.Info.URI)
	}
	description := c.T("music:SongInfo", ms.Queue[0].Track.Info.Title, ms.Queue[0].Track.Info.URI, ms.Queue[0].Track.Info.Author)
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   c.T("music:Duration"),
			Value:  utils.TimeFormat(time.Duration(ms.Queue[0].Track.Info.Length)),
			Inline: true,
		},
		{
			Name:   c.T("music:ETA"),
			Value:  utils.TimeFormat(time.Duration(ms.Queue[0].Track.Info.Length - ms.Player.Position())),
			Inline: true,
		},
	}
	c.Embed.SendEmbed(embed.INFO, description, embed.AddTitle("Now Playing"), embed.AddThumbnail(imgURL), embed.AddProfileFooter(user, c.T), embed.AddFields(fields))
}
