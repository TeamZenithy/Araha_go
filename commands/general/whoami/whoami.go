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
			Usage:                map[string]string{"Required Permission": "**``none``**", "Description": "``Check your permission``", "Usage": fmt.Sprintf("```css\n%swhoami```", utils.Prefix)},
		},
	)
}

const (
	commandName = "whoami"
	commandArg  = "none"
)

func run(ctx handler.CommandContext) error {
	if utils.Contains(utils.Owners, ctx.Message.Author.ID) {
		_, err := ctx.Message.Reply("You are owner")
		return err
	}
	_, err := ctx.Message.Reply("You are not owner")

	return err
}
