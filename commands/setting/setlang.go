package setting

import (
	"github.com/TeamZenithy/Araha/db"
	"github.com/TeamZenithy/Araha/extensions/permissions"
	"github.com/TeamZenithy/Araha/lang"

	"github.com/TeamZenithy/Araha/extensions/embed"
	"github.com/TeamZenithy/Araha/handler"
)

func SetLang() *handler.Cmd {
	return &handler.Cmd{
		Run:      setLang,
		Name:     "setlang",
		Category: "setting",
		Aliases:  []string{""},
		Args:     []string{"lang", "?target"},
		Usage:    "(en | ko | default) --[user | guild]",
	}
}

func setLang(c *handler.Context) {
	langArg, targetArg := c.Arg("lang"), c.Arg("target")

	setGuild := false
	if targetArg == "guild" {
		hasPermission, _ := c.Msg.Member().HasAllPermissions(permissions.MANAGE_GUILD)
		if !hasPermission {
			c.Embed.SendEmbed(embed.BADREQ, c.T("general:BRPermission"))
		}
		setGuild = true
	}

	locale := ""
	if langArg == "default" {
		locale = ""
	} else if lang.IsValidLocale(langArg) {
		locale = langArg
	} else {
		c.Embed.SendEmbed(embed.BADREQ, c.T("locale:BDLocale"))
		return
	}

	if !setGuild {
		if err := db.SetUserLocale(c.Msg.Author.ID, locale); err != nil {
			c.Embed.SendEmbed(embed.ERR_BOT, c.T("error:ErrWriteData"))
			return
		}
		c.Embed.SendEmbed(embed.INFO, c.T("locale:Changed", c.Locale, langArg))
	} else {
		if err := db.SetGuildLocale(c.Msg.GuildID, locale); err != nil {
			c.Embed.SendEmbed(embed.ERR_BOT, c.T("error:ErrWriteData"))
			return
		}
		c.Embed.SendEmbed(embed.INFO, c.T("locale:ChangedServer", c.Locale, langArg))
	}

	return
}
