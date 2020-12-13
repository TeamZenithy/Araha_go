package events

import (
	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	events "github.com/TeamZenithy/Araha/events/audio"
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

//VoiceServerUpdate get voice server change / update events
func VoiceServerUpdate(s *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	logger.Info("Voice Server Update event triggered.")
	vsu := audioengine.VoiceServerUpdate{
		Endpoint: event.Endpoint,
		GuildID:  event.GuildID,
		Token:    event.Token,
	}

	ms, ok := model.Music[event.GuildID]
	if ok && ms.Player != nil {
		if p, err := utils.Lavalink.GetPlayer(event.GuildID); err == nil {
			err = p.Forward(s.State.SessionID, vsu)
			if err != nil {
				logger.Error(err.Error())
			}
			return
		}
	}

	node, err := utils.Lavalink.BestNode()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	handler := events.NewEventHandlerManager()
	utils.Player, err = node.CreatePlayer(event.GuildID, s.State.SessionID, vsu, handler)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	model.Music[event.GuildID].Player = utils.Player
	ms.PlayerCreated <- true
}
