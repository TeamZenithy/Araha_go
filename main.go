package main

import (
	"fmt"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/TeamZenithy/Araha/events"
	"github.com/TeamZenithy/Araha/initializer"
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/sharder"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/TeamZenithy/Araha/web"
)

func main() {
	rawConfig, errFindConfigFile := ioutil.ReadFile("config.toml") // just pass the file name
	if errFindConfigFile != nil {
		logger.Fatal(fmt.Sprintf("Error while load config file: %s", errFindConfigFile.Error()))
		return
	}
	utils.LoadConfig(string(rawConfig))

	initializer.InitLang()

	go web.InitWeb()

	manager := sharder.New("Bot " + utils.Token)
	manager.Name = "Araha"
	manager.LogChannel = utils.ShardLogChannel
	manager.StatusMessageChannel = utils.ShardStatusLogChannel
	// register events
	manager.AddHandler(events.Ready)
	manager.AddHandler(events.MessageCreate)
	manager.AddHandler(events.VoiceServerUpdate)
	manager.AddHandler(events.VoiceStateUpdate)

	recommended, err := manager.GetRecommendedCount()
	if err != nil {
		logger.Fatal("Failed getting recommended shard count")
	}
	if recommended < 2 {
		manager.SetNumShards(5)
	}

	logger.Info("Starting the shard manager")
	manager.Init()
	initializer.InitCommands()

	if err := manager.Start(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start: %s", err.Error()))
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	manager.StopAll()
}
