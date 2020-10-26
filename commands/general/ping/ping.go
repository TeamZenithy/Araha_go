package ping

import (
	"github.com/TeamZenithy/Araha/handler"
)

func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Names:                []string{commandName},
			RequiredArgumentType: []string{commandArg},
			Usage:                map[string]string{"필요한 권한":"**``없음``**", "설명":"``봇이 살아있나 확인합니다.``", "사용법": "```css\n?!ping```"},
		},
	)
}

const (
	commandName = "ping"
	commandArg = "없음"
)

func run(ctx handler.CommandContext) error {
	var _, err = ctx.Message.Reply("Pong!")
	return err
}