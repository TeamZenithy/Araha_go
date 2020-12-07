package events

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/TeamZenithy/Araha/initializer"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

//Ready get discord bot's ready events
func Ready(session *discordgo.Session, event *discordgo.Ready) {
	rawConfig, errFindConfigFile := ioutil.ReadFile("config.toml") // just pass the file name
	if errFindConfigFile != nil {
		fmt.Printf("Error while load config file: %v", errFindConfigFile.Error())
		return
	}
	prefix, errLoadConfigData := utils.GetPrefix(string(rawConfig))
	if errLoadConfigData != nil {
		fmt.Printf("Error while load config data: %v", errLoadConfigData.Error())
	}
	var err = session.UpdateStatus(0, fmt.Sprintf("%shelp", prefix))
	if err != nil {
		log.Println("Error updating status: ", err)
	}
	log.Printf("Logged in as user %s", session.State.User.ID)

	log.Println("Initializing Audio Engine...")

	initializer.InitAudioEngine(event)
}
