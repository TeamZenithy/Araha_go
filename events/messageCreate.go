package events

import (
	"github.com/TeamZenithy/Araha/config"
	"github.com/TeamZenithy/Araha/handler"

	"github.com/bwmarrin/discordgo"
)

//MessageCreate gets message event from discord
func MessageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	go handler.HandleCreatedMessage(session, event, config.Get().Prefix)
}
