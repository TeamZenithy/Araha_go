package queue

import (
	"fmt"

	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"q"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Usage:                map[string]string{"Required Permission": "**``SPEAK``**", "Description": "``Shows current queue``", "Usage": fmt.Sprintf("```css\n%squeue```", utils.Prefix)},
		},
	)
}

const (
	commandName = "queue"
	commandArg  = "none"
)

func run(ctx handler.CommandContext) error {
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil
	}

	var userVoiceState discordgo.VoiceState
	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			userVoiceState = *vs
		}
	}
	if userVoiceState.UserID == "" {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "You are not in a voice channel.")
		return nil
	}
	ms, ok := model.Music[ctx.Message.GuildID]
	if ok && len(ms.Queue) > 0 {
		res := ""
		length := 5
		if len(ms.Queue) < 5 {
			length = 5 - (5 - len(ms.Queue))
		}
		for i, song := range ms.Queue[:length] {
			if i == 0 {
				res += fmt.Sprintf("\n**%d: **%s - (Now Playing)", i+1, song.Track.Info.Title)
			} else {
				res += fmt.Sprintf("\n**%d: **%s", i+1, song.Track.Info.Title)
			}
		}
		if len(ms.Queue) > 5 {
			res += fmt.Sprintf("\n\nAnd **%d** song(s) more in queue", len(ms.Queue)-5)
		}
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Song queue:\n"+res)
	} else {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "There is no music playing in this server.")
	}
	return nil
}
