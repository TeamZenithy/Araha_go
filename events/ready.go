package events

import (
	"fmt"
	"io/ioutil"

	"github.com/TeamZenithy/Araha/initializer"
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

//Ready get discord bot's ready events
func Ready(session *discordgo.Session, event *discordgo.Ready) {
	rawConfig, errFindConfigFile := ioutil.ReadFile("config.toml") // just pass the file name
	if errFindConfigFile != nil {
		logger.Panic(fmt.Sprintf("Error while load config file: %v", errFindConfigFile.Error()))
		return
	}
	prefix, errLoadConfigData := utils.GetPrefix(string(rawConfig))
	if errLoadConfigData != nil {
		logger.Panic(fmt.Sprintf("Error while load config data: %v", errLoadConfigData.Error()))
	}
	var err = session.UpdateStatus(0, fmt.Sprintf("%shelp", prefix))
	if err != nil {
		logger.Warn(fmt.Sprintf("Error updating status: %s", err.Error()))
	}
	logger.Info(fmt.Sprintf("Logged in as user %s#%s(%s)", session.State.User.Username, session.State.User.Discriminator, session.State.User.ID))

	logger.Info("Initializing Audio Engine...")

	initializer.InitAudioEngine(event)
}
