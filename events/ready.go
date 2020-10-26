package events

import (
	"github.com/TeamZenithy/Araha/initializer"
	"github.com/TeamZenithy/Araha/utils"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
)

func Ready(session *discordgo.Session, event *discordgo.Ready) {
	rawConfig, errFindConfigFile := ioutil.ReadFile("config.toml") // just pass the file name
	if errFindConfigFile != nil {
		fmt.Println("Error while load config file: " + errFindConfigFile.Error())
		return
	}
	errLoadConfigData, prefix := utils.GetPrefix(string(rawConfig))
	if  errLoadConfigData != nil {
		fmt.Println("Error while load config data: " + errLoadConfigData.Error())
	}
	var err = session.UpdateStatus(0, prefix + "help")
	if err != nil {
		log.Println("Error updating status: ", err)
	}
	log.Println("Logged in as user " + session.State.User.ID)

	log.Println("Initializing Audio Engine...")

	initializer.InitAudioEngine(event)

}
