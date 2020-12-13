package events

import (
	"fmt"
	"io/ioutil"

	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/utils"

	"github.com/bwmarrin/discordgo"
)

//MessageCreate gets message event from discord
func MessageCreate(session *discordgo.Session, event *discordgo.MessageCreate) {
	rawConfig, errFindConfigFile := ioutil.ReadFile("config.toml") // just pass the file name
	if errFindConfigFile != nil {
		logger.Error(fmt.Sprintf("Error while load config file: %v", errFindConfigFile.Error()))
		return
	}
	prefix, errLoadConfigData := utils.GetPrefix(string(rawConfig))
	if errLoadConfigData != nil {
		logger.Error(fmt.Sprintf("Error while load config data: %v", errLoadConfigData.Error()))
	}
	go handler.HandleCreatedMessage(session, event, prefix)
}
