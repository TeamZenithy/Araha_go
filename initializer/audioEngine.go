package initializer

import (
	"github.com/TeamZenithy/Araha/config"
	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

//InitAudioEngine initialize lavalink client
func InitAudioEngine(event *discordgo.Ready) {
	utils.Lavalink = audioengine.NewLavalink("1", event.User.ID)
	host, port, pass := config.Get().LavalinkHost, config.Get().LavalinkPort, config.Get().LavalinkPass

	err := utils.Lavalink.AddNodes(audioengine.NodeConfig{
		REST:      "http://" + host + ":" + port,
		WebSocket: "ws://" + host + ":" + port,
		Password:  pass,
	})
	model.Music = make(map[string]*model.MusicStruct)

	if err != nil {
		logger.Error(err.Error())
	}
}
