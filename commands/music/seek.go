package music

import (
	"strconv"

	"github.com/TeamZenithy/Araha/extensions/embed"
	m "github.com/TeamZenithy/Araha/middlewares"
	"github.com/TeamZenithy/Araha/model"

	"github.com/TeamZenithy/Araha/handler"
)

func Seek() *handler.Cmd {
	return &handler.Cmd{
		Run:         seek,
		Name:        "seek",
		Category:    "music",
		Aliases:     []string{""},
		Args:        []string{"time"},
		Middlewares: []handler.HandlerFunc{m.LoadQueue(), m.VoiceWithMusic()},
		Usage:       "(second)",
	}
}

func seek(c *handler.Context) {
	ms := c.Get("queue").(*model.MusicStruct)

	pos, errNotSecond := strconv.Atoi(c.Arg("time"))
	pos = pos * 1000
	if errNotSecond != nil || pos < 0 {
		c.Embed.SendEmbed(embed.BADREQ, c.T("music:BRTime"))
		return
	}
	if ms.Queue[0].Track.Info.Length <= int64(pos) {
		c.Embed.SendEmbed(embed.BADREQ, c.T("music:BRTime2"))
		return
	}
	ms.Player.Seek(int64(pos))
	c.Embed.SendEmbed(embed.BADREQ, c.T("music:SeekTo", c.Arg("time")))
}
