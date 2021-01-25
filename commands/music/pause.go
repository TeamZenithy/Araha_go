package music

import (
	"github.com/TeamZenithy/Araha/extensions/embed"
	m "github.com/TeamZenithy/Araha/middlewares"
	"github.com/TeamZenithy/Araha/model"

	"github.com/TeamZenithy/Araha/handler"
)

func Pause() *handler.Cmd {
	return &handler.Cmd{
		Run:         pause,
		Name:        "pause",
		Category:    "music",
		Aliases:     []string{""},
		Args:        []string{},
		Middlewares: []handler.HandlerFunc{m.LoadQueue(), m.VoiceWithMusic()},
		Usage:       "",
	}
}

func pause(c *handler.Context) {
	queue := c.Get("queue").(*model.MusicStruct)
	if queue.Player.Paused() {
		c.Embed.SendEmbed(embed.BADREQ, c.T("music:AlreadyPaused"))
	} else {
		queue.Player.Pause(true)
		c.Embed.SendEmbed(embed.BADREQ, c.T("music:Paused"))
	}
	return
}
