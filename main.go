package main

import (
	"github.com/TeamZenithy/Araha/events"
	"github.com/TeamZenithy/Araha/initializer"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rawConfig, errFindConfigFile := ioutil.ReadFile("config.toml") // just pass the file name
	if errFindConfigFile != nil {
		log.Fatalln("Error while load config file: " + errFindConfigFile.Error())
		return
	}
	errLoadConfigData, token := utils.GetToken(string(rawConfig))
	if errLoadConfigData != nil {
		log.Fatalln("Error while load config data: " + errLoadConfigData.Error())
	}

	var bot, err = discordgo.New("Bot " + token)
	// register events
	bot.AddHandler(events.Ready)
	bot.AddHandler(events.MessageCreate)
	bot.AddHandler(events.VoiceServerUpdate)

	initializer.InitCommands()
	err = bot.Open()

	if err != nil {
		log.Fatalln("Error opening Discord session: ", err)
	}

	log.Println("Bot is now running.")

	// wait forever
	select {}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Close()
}
