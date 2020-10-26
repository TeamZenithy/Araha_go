package join

import (
	"github.com/TeamZenithy/Araha/handler"
	"log"
)

func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Names:                []string{commandName},
			RequiredArgumentType: []string{commandArg},
			Usage:                map[string]string{"필요한 권한":"**``음성 채널 접속``**", "설명":"``유저가 접속해있는 음성채널에 접속합니다.``", "사용법": "```css\n?!join```"},
		},
	)
}

const (
	commandName ="join"
	commandArg = "없음"
)

func run(ctx handler.CommandContext) error {
	c, err := ctx.Session.State.Channel(ctx.Message.ChannelID)
	if err != nil {
		log.Println("fail find channel")
		ctx.Message.Reply("You are not in voice channel!")
		return nil
	}

	g, err := ctx.Session.State.Guild(c.GuildID)
	if err != nil {
		log.Println("fail find guild")
		ctx.Message.Reply("You are not in guild!")
		return nil
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			log.Println("trying to connect to channel")
			err = ctx.Session.ChannelVoiceJoinManual(c.GuildID, vs.ChannelID, false, false)
			if err != nil {
				log.Println(err)
				ctx.Message.Reply("Failed to connect to voice channel!")
			} else {
				log.Println("channel voice join succeeded")
			}
		}
	}
	ctx.Message.Reply("음성채널에 접속했습니다!")
	return nil
}
