package locale

import (
	"fmt"
	"strings"

	"github.com/TeamZenithy/Araha/db"

	"github.com/TeamZenithy/Araha/extensions/embed"
	"github.com/TeamZenithy/Araha/extensions/permissions"
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
			Category:             utils.CATEGORY_LOCALE,
			Description:          &handler.Description{ReqPermsission: "none", Usage: "setlang [user | guild] (en | ko | default)"},
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
	if query == "" {
		e.SendEmbed(embed.BADREQ, ctx.T("locale:BDEnter"))
		return nil
	}

	setGuild := false
	if strings.Contains(query, "guild") {
		hasPermission, _ := utils.MemberHasPermission(ctx.Session, ctx.Message.GuildID, ctx.Message.Author.ID, permissions.MANAGE_GUILD)
		if !hasPermission {
			e.SendEmbed(embed.BADREQ, ctx.T("general:BRPermission"))
		}
		setGuild = true
		query = strings.TrimSpace(strings.ReplaceAll(query, "guild", ""))
	}

	locale := ""
	if query == "default" {
		locale = ""
	} else if lang.IsValidLocale(query) {
		locale = query
	} else {
		fmt.Println()
		e.SendEmbed(embed.BADREQ, ctx.T("locale:BDLocale"))
		return nil
	}

	if !setGuild {
		if err := db.SetUserLocale(ctx.Message.Author.ID, locale); err != nil {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("error:ErrWriteData"))
			return nil
		}
		e.SendEmbed(embed.INFO, ctx.T("locale:Changed", ctx.Locale, query))
	} else {
		if err := db.SetGuildLocale(ctx.Message.GuildID, locale); err != nil {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("error:ErrWriteData"))
			return nil
		}
		e.SendEmbed(embed.INFO, ctx.T("locale:ChangedServer", ctx.Locale, query))
	}

	return nil
}
