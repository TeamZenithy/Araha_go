package whoami

import (
	"fmt"

	"github.com/TeamZenithy/Araha/utils"

	"github.com/TeamZenithy/Araha/handler"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"wai"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_GENERAL,
			Usage:                map[string]string{"필요한 권한": "**``없음``**", "설명": "``자신이 누구인지 확인합니다.``", "사용법": fmt.Sprintf("```css\n%swhoami```", utils.Prefix)},
		},
	)
}

const (
	commandName = "whoami"
	commandArg  = "없음"
)

func run(ctx handler.CommandContext) error {
	if utils.Contains(utils.Owners, ctx.Message.Author.ID) {
		_, err := ctx.Message.Reply("You are owner")
		return err
	}
	_, err := ctx.Message.Reply("You are not owner")

	return err
}
