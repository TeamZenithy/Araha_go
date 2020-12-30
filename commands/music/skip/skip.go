package skip

import (
	"math"
	"strconv"

	"github.com/TeamZenithy/Araha/extensions/embed"
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
			Description:          &handler.Description{ReqPermsission: "SPEAK", Usage: "skip"},
		},
	)
}

const (
	commandName = "skip"
	commandArg  = "none"
)

func run(ctx handler.CommandContext) error {
	e := embed.New(ctx.Session, ctx.Message.ChannelID)

	if ms, ok := model.Music[ctx.Message.GuildID]; ok {
		guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
		if err != nil {
			return nil
		}

		if isInVoice := utils.IsInVoiceWithMusic(guild, ctx.Message.Author.ID); !isInVoice {
			e.SendEmbed(embed.BADREQ, ctx.T("music:BRNotPlaying"))
			return nil
		}

		usersInVoice := math.Floor(float64(utils.GetUsersInVoice(guild) / 2))
		skips := ms.Queue[0].Skips
		requirement := (skips + 1) / float64(usersInVoice)
		if usersInVoice <= 2 || requirement >= 0.4 {
			ms.Player.Stop()
		} else {
			skips++
			e.SendEmbed(embed.INFO, ctx.T("music:VoteAdded", strconv.Itoa(int(usersInVoice-skips)), strconv.Itoa(int(skips)), strconv.Itoa(int(usersInVoice))))
		}
	} else {
		e.SendEmbed(embed.BADREQ, ctx.T("music:NoMusic"))
	}
	return nil
}
