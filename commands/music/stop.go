package music

import (
	"github.com/TeamZenithy/Araha/extensions/embed"
	m "github.com/TeamZenithy/Araha/middlewares"
	"github.com/TeamZenithy/Araha/utils"

	"github.com/TeamZenithy/Araha/handler"
)

func Stop() *handler.Cmd {
	return &handler.Cmd{
		Run:         stop,
		Name:        "stop",
		Category:    "music",
		Aliases:     []string{""},
		Args:        []string{},
		Middlewares: []handler.HandlerFunc{m.LoadQueue(), m.VoiceWithMusic()},
		Usage:       "",
	}
}

func stop(c *handler.Context) {
	if returnedMessage := utils.LeaveAndDestroy(c.Session, c.Msg.GuildID); returnedMessage != "" {
		c.Session.ChannelMessageSend(c.Msg.ChannelID, c.T("music:ErrStopped")+"\n"+returnedMessage)
	} else {
		c.Embed.SendEmbed(embed.INFO, c.T("music:Stopped"))
	}
}
