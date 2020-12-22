package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/TeamZenithy/Araha/events"
	"github.com/TeamZenithy/Araha/initializer"
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

func main() {
	rawConfig, errFindConfigFile := ioutil.ReadFile("config.toml") // just pass the file name
	if errFindConfigFile != nil {
		logger.Fatal(fmt.Sprintf("Error while load config file: %s", errFindConfigFile.Error()))
		return
	}
	utils.LoadConfig(string(rawConfig))

	var bot, err = discordgo.New("Bot " + utils.Token)
	// register events
	bot.AddHandler(events.Ready)
	bot.AddHandler(events.MessageCreate)
	bot.AddHandler(events.VoiceServerUpdate)
	bot.AddHandler(events.VoiceStateUpdate)

	initializer.InitCommands()
	logger.Info("Trying to log in...")
	err = bot.Open()

	if err != nil {
		logger.Fatal(fmt.Sprintf("Error opening Discord session: %v", err))
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Close()
}
