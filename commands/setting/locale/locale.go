package locale

import (
	"fmt"

	"github.com/TeamZenithy/Araha/db"

	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/lang"
	"github.com/TeamZenithy/Araha/utils"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"locale", "setlocale"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_GENERAL,
			Usage:                map[string]string{"Required Permission": "**``none``**", "Description": "``Set User or Guild's Language``", "Usage": fmt.Sprintf("```css\n%ssetlang [--user | --guild] (en | ko)```", utils.Prefix)},
		},
	)
}

const (
	commandName = "setlang"
	commandArg  = "locale"
)

func run(ctx handler.CommandContext) error {
	query := ctx.Arguments[commandArg]
	if query == "" {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Please provide locale.")
		return nil
	}

	if !lang.IsValidLocale(query) {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Enter Valid locale.")
		return nil
	}

	err := db.SetUserLocale(ctx.Message.Author.ID, query)
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error while writing data")
		return nil
	}

	ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Successful!")

	return nil
}
