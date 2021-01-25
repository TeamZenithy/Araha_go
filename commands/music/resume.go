package music

import (
	"github.com/TeamZenithy/Araha/extensions/embed"
	m "github.com/TeamZenithy/Araha/middlewares"
	"github.com/TeamZenithy/Araha/model"

	"github.com/TeamZenithy/Araha/handler"
)

func Resume() *handler.Cmd {
	return &handler.Cmd{
		Run:         resume,
		Name:        "resume",
		Category:    "music",
		Aliases:     []string{""},
		Args:        []string{},
		Middlewares: []handler.HandlerFunc{m.LoadQueue(), m.VoiceWithMusic()},
		Usage:       "",
	}
}

func resume(c *handler.Context) {
	queue := c.Get("queue").(*model.MusicStruct)
	if !queue.Player.Paused() {
		c.Embed.SendEmbed(embed.BADREQ, c.T("music:AlreadyPlaying"))
	} else {
		queue.Player.Pause(false)
		c.Embed.SendEmbed(embed.BADREQ, c.T("music:Resume"))
	}
	return
}
