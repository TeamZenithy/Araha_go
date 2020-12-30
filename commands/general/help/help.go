package help

import (
	"strings"

	"github.com/TeamZenithy/Araha/extensions/embed"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"

	"github.com/TeamZenithy/Araha/handler"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"h", "guide", "manual"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_GENERAL,
			Description:          &handler.Description{ReqPermsission: "none", Usage: "help [command]"},
		},
	)
}

const (
	commandName = "help"
	commandArg  = "Command"
)

func run(ctx handler.CommandContext) error {
	e := embed.New(ctx.Session, ctx.Message.ChannelID)
	if _, exists := ctx.Arguments[commandArg]; exists {
		// show specific information about a command
		cmd, isExists := handler.Commands[strings.ToLower(ctx.Arguments[commandArg])]
		aliasCommand, isExistsAliasCommand := handler.Aliases[strings.ToLower(ctx.Arguments[commandArg])]

		if !isExists && isExistsAliasCommand {
			isExists = isExistsAliasCommand
			cmd = handler.Commands[aliasCommand]
		}

		if isExists {
			fields := []*discordgo.MessageEmbedField{}
			alias := ""
			if len(cmd.Aliases) < 1 {
				alias = ctx.T("generalNone")
			}
			for i, d := range cmd.Aliases {
				if i != 0 {
					alias += ", "
				}
				alias += d
			}

			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   ctx.T("general:Alias"),
				Value:  alias,
				Inline: true,
			})

			reqPermission := cmd.Description.ReqPermsission
			if reqPermission == "none" {
				reqPermission = ctx.T("general:None")
			}

			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   ctx.T("general:Permission"),
				Value:  reqPermission,
				Inline: true,
			})

			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   ctx.T("general:Usage"),
				Value:  cmd.Description.Usage,
				Inline: true,
			})

			e.SendEmbed(embed.INFO, ctx.T("cmd:"+cmd.Name), embed.AddTitle(cmd.Name), embed.AddFields(fields))

			return nil
		}
		e.SendEmbed(embed.BADREQ, ctx.T("general:BRNotFound"))
		return nil
	}

	list := map[string]string{}
	fields := []*discordgo.MessageEmbedField{}
	for _, v := range handler.Commands {
		target := ""
		switch v.Category {
		case utils.CATEGORY_GENERAL:
			target = "General"
		case utils.CATEGORY_LOCALE:
			target = "Locale"
		case utils.CATEGORY_MUSIC:
			target = "Music"
		}
		list[target] += v.Name + " "
	}

	for k, v := range list {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   k,
			Value:  strings.ReplaceAll(strings.TrimSpace(v), " ", ", "),
			Inline: false,
		})
	}
	e.SendEmbed(embed.INFO, "", embed.AddTitle(ctx.T("general:CommandList")), embed.AddFields(fields))
	return nil
}
