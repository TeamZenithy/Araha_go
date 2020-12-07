package events

import (
	"log"

	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

func VoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	ms, ok := model.Music[event.GuildID]
	if !ok {
		return
	}

	guild, err := session.State.Guild(event.GuildID)
	if err != nil {
		log.Printf("No guild found in State for %s: %s", event.GuildID, err)
		return
	}
	if utils.GetUsersInVoice(guild) == 0 {
		ms.SongEnd <- "end"
		if returnedMessage := utils.LeaveAndDestroy(session, event.GuildID); returnedMessage != "" {
			log.Print(returnedMessage)
		}
		return
	}
}