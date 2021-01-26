package music

import (
	"math/rand"
	"time"

	"github.com/TeamZenithy/Araha/extensions/embed"
	m "github.com/TeamZenithy/Araha/middlewares"
	"github.com/TeamZenithy/Araha/model"

	"github.com/TeamZenithy/Araha/handler"
)

func Shuffle() *handler.Cmd {
	return &handler.Cmd{
		Run:         shuffle,
		Name:        "shuffle",
		Category:    "music",
		Aliases:     []string{""},
		Args:        []string{},
		Middlewares: []handler.HandlerFunc{m.LoadQueue(), m.VoiceWithMusic()},
		Usage:       "",
	}
}

func shuffle(c *handler.Context) {
	ms := c.Get("queue").(*model.MusicStruct)

	queue := ms.Queue[1:]
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(queue), func(i, j int) {
		queue[i], queue[j] = queue[j], queue[i]
	})
	queue = append([]model.Song{ms.Queue[0]}, queue...)
	ms.Queue = queue
	c.Embed.SendEmbed(embed.INFO, c.T("music:Shuffled"))
}
