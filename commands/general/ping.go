package general

import (
	"strconv"

	"github.com/TeamZenithy/Araha/extensions/embed"
	"github.com/TeamZenithy/Araha/handler"
)

func Ping() *handler.Cmd {
	return &handler.Cmd{
		Run:      ping,
		Name:     "ping",
		Category: "general",
		Aliases:  []string{"pong"},
		Args:     []string{},
		Usage:    "",
	}
}

func ping(c *handler.Context) {
	c.Embed.SendEmbed(embed.INFO, "Pong! "+strconv.Itoa(int(c.Session.HeartbeatLatency().Milliseconds()))+"ms :stopwatch:")
	return
}
