package stop

import (
	"fmt"

	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/utils"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"quit"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Usage:                map[string]string{"필요한 권한": "**``음성 채널 발언권``**", "설명": "``노래를 중지하고 대기열을 초기화 합니다.``", "사용법": fmt.Sprintf("```css\n%sstop```", utils.Prefix)},
		},
	)
}

const (
	commandName = "stop"
	commandArg  = "없음"
)

func run(ctx handler.CommandContext) error {
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil
	}

	if utils.IsInVoiceWithMusic(guild, ctx.Message.Author.ID) {
		if returnedMessage := utils.LeaveAndDestroy(ctx.Session, ctx.Message.GuildID); returnedMessage != "" {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Playback may not have been stopped.\n"+returnedMessage)
		} else {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Playback stopped.")
		}
	} else {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "You're not listening to my music :(")
	}
	return nil
}
