package main

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/TeamZenithy/Araha/config"
	"github.com/TeamZenithy/Araha/db"
	"github.com/TeamZenithy/Araha/events"
	"github.com/TeamZenithy/Araha/initializer"
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/sharder"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/TeamZenithy/Araha/web"
)

func main() {
	config.LoadConfig()

	initializer.InitLang()
	db.InitRedis()

	go web.InitWeb()

	manager := sharder.New("Bot " + config.Get().Token)
	manager.Name = "Araha"
	manager.LogChannel = config.Get().ShardLogChannel
	manager.StatusMessageChannel = config.Get().ShardStatusLogChannel
	// register events
	manager.AddHandler(events.Ready)
	manager.AddHandler(events.MessageCreate)
	manager.AddHandler(events.VoiceServerUpdate)
	manager.AddHandler(events.VoiceStateUpdate)

	if config.Get().Release {
		recommended, err := manager.GetRecommendedCount()
		if err != nil {
			logger.Fatal("Failed getting recommended shard count")
		}
		if recommended < 2 {
			manager.SetNumShards(5)
		}
	} else {
		manager.SetNumShards(1)
	}

	logger.Info("Starting the shard manager")
	manager.Init()
	initializer.InitCommands()

	if err := manager.Start(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start: %s", err.Error()))
	}

	sharder.CurManager = manager

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	logger.Info("Saving all data to db before shutting down...")
	utils.RDB.Save(context.Background())
	manager.StopAll()
	logger.Info("All processes have been terminated.")
}
