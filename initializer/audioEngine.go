package initializer

import (
	"io/ioutil"
	"log"

	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

//InitAudioEngine initialize lavalink client
func InitAudioEngine(event *discordgo.Ready) {
	rawConfig, errFindConfigFile := ioutil.ReadFile("config.toml") // just pass the file name
	if errFindConfigFile != nil {
		log.Fatalln("Error while load config file: " + errFindConfigFile.Error())
		return
	}
	utils.Lavalink = audioengine.NewLavalink("1", event.User.ID)
	host, port, pass, errLoadConfigData := utils.GetLavalinkConfig(string(rawConfig))
	if errLoadConfigData != nil {
		log.Fatalln("Error while load config data: " + errLoadConfigData.Error())
	}

	err := utils.Lavalink.AddNodes(audioengine.NodeConfig{
		REST:      "http://" + host + ":" + port,
		WebSocket: "ws://" + host + ":" + port,
		Password:  pass,
	})
	model.Music = make(map[string]*model.MusicStruct)

	if err != nil {
		log.Println(err)
	}
}
