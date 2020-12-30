package ping

import (
	"strconv"

	"github.com/TeamZenithy/Araha/extensions/embed"
	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/utils"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"pong"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_GENERAL,
			Description:          &handler.Description{ReqPermsission: "none", Usage: "ping"},
		},
	)
}

const (
	commandName = "ping"
	commandArg  = "none"
)

func run(ctx handler.CommandContext) error {
	e := embed.New(ctx.Session, ctx.Message.ChannelID)
	e.SendEmbed(embed.INFO, "Pong! "+strconv.Itoa(int(ctx.Session.HeartbeatLatency().Milliseconds()))+"ms :stopwatch:")
	return nil
}
