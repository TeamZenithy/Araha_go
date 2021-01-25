package music

import (
	m "github.com/TeamZenithy/Araha/middlewares"
	"github.com/TeamZenithy/Araha/model"

	"github.com/TeamZenithy/Araha/handler"
)

func Skip() *handler.Cmd {
	return &handler.Cmd{
		Run:         skip,
		Name:        "skip",
		Category:    "music",
		Aliases:     []string{"s"},
		Args:        []string{},
		Middlewares: []handler.HandlerFunc{m.LoadQueue(), m.VoiceWithMusic()},
		Usage:       "",
	}
}

func skip(c *handler.Context) {
	ms := c.Get("queue").(*model.MusicStruct)

	ms.Player.Stop()
}
