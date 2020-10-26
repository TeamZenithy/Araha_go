package events

import (
	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
	"log"
)

func VoiceServerUpdate(s *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	log.Println("received Voice Server Update event.")
	vsu := audioengine.VoiceServerUpdate{
		Endpoint: event.Endpoint,
		GuildID:  event.GuildID,
		Token:    event.Token,
	}

	if p, err := utils.Lavalink.GetPlayer(event.GuildID); err == nil {
		err = p.Forward(s.State.SessionID, vsu)
		if err != nil {
			log.Println(err)
		}
		return
	}

	node, err := utils.Lavalink.BestNode()
	if err != nil {
		log.Println(err)
		return
	}

	handler := new(audioengine.DummyEventHandler)
	utils.Player, err = node.CreatePlayer(event.GuildID, s.State.SessionID, vsu, handler)
	if err != nil {
		log.Println(err)
		return
	}
}
