package shuffle

import (
	"fmt"
	"math/rand"
	"time"

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
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Usage:                map[string]string{"Required Permission": "**``SPEAK``**", "Description": "``Shuffle the queue``", "Usage": fmt.Sprintf("```css\n%sshuffle```", utils.Prefix)},
		},
	)
}

const (
	commandName = "shuffle"
	commandArg  = "none"
)

func run(ctx handler.CommandContext) error {
	if ms, ok := model.Music[ctx.Message.GuildID]; ok {
		guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
		if err != nil {
			return nil
		}

		if isInVoice := utils.IsInVoiceWithMusic(guild, ctx.Message.Author.ID); !isInVoice {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:BRNotPlaying"))
			return nil
		}

		queue := ms.Queue[1:]
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(queue), func(i, j int) {
			queue[i], queue[j] = queue[j], queue[i]
		})
		queue = append([]model.Song{ms.Queue[0]}, queue...)
		ms.Queue = queue
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:Shuffled"))
	} else {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("music:NoMusic"))
	}
	return nil
}
