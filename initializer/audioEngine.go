package initializer

import (
	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
	"log"
)

func InitAudioEngine(event *discordgo.Ready) {
	utils.Lavalink = audioengine.NewLavalink("1", event.User.ID)
	err := utils.Lavalink.AddNodes(audioengine.NodeConfig{
		REST: "http://localhost:2333",
		WebSocket: "ws://localhost:2333",
		Password: "7262rcrw3392",
	})
	if err != nil {
		log.Println(err)
	}
}
