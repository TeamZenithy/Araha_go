package ping

import (
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
			Usage:                map[string]string{"필요한 권한": "**``없음``**", "설명": "``봇이 살아있나 확인합니다.``", "사용법": "```css\n?!ping```"},
		},
	)
}

const (
	commandName = "ping"
	commandArg  = "없음"
)

func run(ctx handler.CommandContext) error {
	var _, err = ctx.Message.Reply("Pong!")
	return err
}
