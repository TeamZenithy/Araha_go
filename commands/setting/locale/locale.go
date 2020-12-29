package locale

import (
	"fmt"
	"strings"

	"github.com/TeamZenithy/Araha/db"

	"github.com/TeamZenithy/Araha/extensions/embed"
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
			Aliases:              []string{"locale"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_GENERAL,
			Usage:                map[string]string{"Required Permission": "**``none``**", "Description": "``Set User or Guild's Language``", "Usage": fmt.Sprintf("```css\n%ssetlang [--user | --guild] (en | ko | default)```", utils.Prefix)},
		},
	)
}

const (
	commandName = "setlang"
	commandArg  = "locale"
)

func run(ctx handler.CommandContext) error {
	e := embed.New(ctx.Session, ctx.Message.ChannelID)

	query := ctx.Arguments[commandArg]
	fmt.Println(query)
	if query == "" {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Please provide locale.")
		return nil
	}

	setGuild := false
	if strings.Contains(query, "--guild") {
		// TODO: Check Permission
		if false {
			e.SendEmbed(embed.BADREQ, ctx.T("common:Permission"))
		}
		setGuild = true
	}

	locale := query
	if query == "default" {
		locale = ""
	}

	if !lang.IsValidLocale(query) && locale != "" {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("locale:BDLocale"))
		return nil
	}

	if !setGuild {
		if err := db.SetUserLocale(ctx.Message.Author.ID, locale); err != nil {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error while writing data")
			return nil
		}
		e.SendEmbed(embed.INFO, ctx.T("locale:Changed", locale))
	} else {
		if err := db.SetGuildLocale(ctx.Message.GuildID, locale); err != nil {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error while writing data")
			return nil
		}
		e.SendEmbed(embed.INFO, ctx.T("locale:ChangedServer", locale))
	}

	return nil
}
