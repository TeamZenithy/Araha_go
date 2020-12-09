package queue

import (
	"fmt"

	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"q"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Usage:                map[string]string{"필요한 권한": "**``음성 채널 발언권``**", "설명": "``재생 대기열을 보여줍니다.``", "사용법": "```css\n?!queue```"},
		},
	)
}

const (
	commandName = "queue"
	commandArg  = "없음"
)

func run(ctx handler.CommandContext) error {
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil
	}

	var userVoiceState discordgo.VoiceState
	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			userVoiceState = *vs
		}
	}
	if userVoiceState.UserID == "" {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "You are not in a voice channel.")
		return nil
	}
	ms, ok := model.Music[ctx.Message.GuildID]
	if ok && len(ms.Queue) > 0 {
		res := ""
		for i, song := range ms.Queue {
			res += fmt.Sprintf("\n**%d: **%s", i+1, song.Track.Info.Title)
		}
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Song queue:\n"+res)
	} else {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "There is no music playing in this server.")
	}
	return nil
}
