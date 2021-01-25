package music

import (
	"strconv"

	"github.com/TeamZenithy/Araha/extensions/embed"
	"github.com/TeamZenithy/Araha/logger"
	m "github.com/TeamZenithy/Araha/middlewares"
	"github.com/TeamZenithy/Araha/model"

	"github.com/TeamZenithy/Araha/handler"
)

func Volume() *handler.Cmd {
	return &handler.Cmd{
		Run:         volume,
		Name:        "volume",
		Category:    "music",
		Aliases:     []string{"vol"},
		Args:        []string{"?volume"},
		Middlewares: []handler.HandlerFunc{m.LoadQueue(), m.VoiceWithMusic()},
		Usage:       "[0~200]",
	}
}

func volume(c *handler.Context) {
	ms := c.Get("queue").(*model.MusicStruct)
	volumeArg := c.Arg("volume")
	if volumeArg == "" {
		c.Embed.SendEmbed(embed.INFO, c.T("music:CurVolume", strconv.Itoa(ms.Player.GetVolume())))
		return
	}
	vol, err := strconv.Atoi(volumeArg)
	if err != nil {
		c.Embed.SendEmbed(embed.BADREQ, c.T("music:BRNumber"))
		return
	}
	if vol < 1 || vol > 200 {
		c.Embed.SendEmbed(embed.BADREQ, c.T("music:BRNumber"))
		return
	}
	prevVol := ms.Player.GetVolume()
	err = ms.Player.Volume(vol)
	if err != nil {
		logger.Warn(err.Error())
	}
	c.Embed.SendEmbed(embed.INFO, c.T("music:Volume", strconv.Itoa(prevVol), strconv.Itoa(ms.Player.GetVolume())))
	return
}
