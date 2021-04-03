package general

import (
	"strings"

	"github.com/TeamZenithy/Araha/config"
	"github.com/TeamZenithy/Araha/extensions/embed"
	"github.com/TeamZenithy/Araha/handler"
	"github.com/bwmarrin/discordgo"
)

func Help() *handler.Cmd {
	return &handler.Cmd{
		Run:      help,
		Name:     "help",
		Category: "general",
		Aliases:  []string{"h", "guide", "manual"},
		Args:     []string{"?cmdName"},
		Usage:    "[command]",
	}
}

func help(c *handler.Context) {
	arg := c.Arg("cmdName")
	if arg != "" {
		cmd := handler.FindCommand(c.Arg("cmdName"))

		if cmd != nil {
			fields := []*discordgo.MessageEmbedField{}
			alias := ""
			if len(cmd.Aliases) < 1 {
				alias = c.T("generalNone")
			}
			for i, d := range cmd.Aliases {
				if i != 0 {
					alias += ", "
				}
				alias += d
			}

			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   c.T("general:Alias"),
				Value:  alias,
				Inline: true,
			})

			// reqPermission := cmd.Description.ReqPermsission
			// if reqPermission == "none" {
			// 	reqPermission = c.T("general:None")
			// }

			// fields = append(fields, &discordgo.MessageEmbedField{
			// 	Name:   c.T("general:Permission"),
			// 	Value:  reqPermission,
			//		Inline: true,
			//		})

			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   c.T("general:Usage"),
				Value:  config.Get().Prefix + cmd.Name + " " + cmd.Usage,
				Inline: true,
			})

			c.Embed.SendEmbed(embed.INFO, c.T("cmd:"+cmd.Name), embed.AddTitle(cmd.Name), embed.AddFields(fields))
		} else {
			c.Embed.SendEmbed(embed.BADREQ, c.T("general:BRNotFound"))
		}
		return
	}

	list := map[string]string{}
	fields := []*discordgo.MessageEmbedField{}
	for _, v := range handler.Commands {
		list[v.Category] += v.Name + " "
	}

	for k, v := range list {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   k,
			Value:  strings.ReplaceAll(strings.TrimSpace(v), " ", ", "),
			Inline: false,
		})
	}
	c.Embed.SendEmbed(embed.INFO, "", embed.AddTitle(c.T("general:CommandList")), embed.AddFields(fields))
	return
}
