package events

import (
	"fmt"

	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

//VoiceStateUpdate handles voice state update event
func VoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	ms, ok := model.Music[event.GuildID]
	if !ok {
		return
	}

	guild, err := session.State.Guild(event.GuildID)
	if err != nil {
		logger.Warn(fmt.Sprintf("No guild found in State for %s: %s", event.GuildID, err))
		return
	}
	if utils.GetUsersInVoice(guild) == 0 {
		ms.SongEnd <- "end"
		if returnedMessage := utils.LeaveAndDestroy(session, event.GuildID); returnedMessage != "" {
			logger.Info(returnedMessage)
		}
		return
	}
}
