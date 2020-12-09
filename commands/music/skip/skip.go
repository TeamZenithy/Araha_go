package skip

import (
	"fmt"
	"math"

	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"s"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Usage:                map[string]string{"필요한 권한": "**``음성 채널 발언권``**", "설명": "``현재 재생중인 노래를 건너뜁니다.``", "사용법": "```css\n?!skip```"},
		},
	)
}

const (
	commandName = "skip"
	commandArg  = "없음"
)

func run(ctx handler.CommandContext) error {
	if ms, ok := model.Music[ctx.Message.GuildID]; ok {
		guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
		if err != nil {
			return nil
		}

		if isInVoice := utils.IsInVoiceWithMusic(guild, ctx.Message.Author.ID); !isInVoice {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "You're not listening to my music :(")
			return nil
		}

		usersInVoice := math.Floor(float64(utils.GetUsersInVoice(guild) / 2))
		skips := ms.Queue[0].Skips
		requirement := (skips + 1) / float64(usersInVoice)
		if usersInVoice <= 2 || requirement >= 0.4 {
			ms.Player.Stop()
		} else {
			skips++
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("Vote added! Need %d more (%d/%d).", int(usersInVoice-skips), int(skips), int(usersInVoice)))
		}
	} else {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "There is no music playing.")
	}
	return nil
}
