package middlewares

import (
	"github.com/TeamZenithy/Araha/extensions/embed"
	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/model"
	"github.com/bwmarrin/discordgo"
)

func UserVoiceState() handler.HandlerFunc {
	return func(c *handler.Context) bool {
		guild := c.Msg.Guild()
		var userVoiceState discordgo.VoiceState
		for _, vs := range guild.VoiceStates {
			if vs.UserID == c.Msg.Author.ID {
				userVoiceState = *vs
			}
		}
		if userVoiceState.UserID == "" {
			c.Embed.SendEmbed(embed.BADREQ, c.T("music:NotInVoiceChannel"))
			return false
		}

		c.Set("userVoiceState", &userVoiceState)
		return true
	}
}

func VoiceWithMusic() handler.HandlerFunc {
	return func(c *handler.Context) bool {
		ms := c.Get("queue").(*model.MusicStruct)
		guild := c.Msg.Guild()
		for _, vs := range guild.VoiceStates {
			if ms.ChannelID == vs.ChannelID && c.Msg.Author.ID == vs.UserID {
				return true
			}
		}
		c.Embed.SendEmbed(embed.BADREQ, c.T("music:NoMusic"))
		return false
	}

}

func LoadQueue() handler.HandlerFunc {
	return func(c *handler.Context) bool {
		ms, ok := model.Music[c.Msg.GuildID]
		if !ok || len(ms.Queue) < 1 {
			c.Embed.SendEmbed(embed.BADREQ, c.T("music:NoMusic"))
			return false
		}
		c.Set("queue", ms)
		return true
	}
}
