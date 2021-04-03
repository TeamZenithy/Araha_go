package events

import (
	"fmt"
	"strings"

	"github.com/TeamZenithy/Araha/config"
	"github.com/TeamZenithy/Araha/initializer"
	"github.com/TeamZenithy/Araha/logger"
	"github.com/bwmarrin/discordgo"
)

//Ready get discord bot's ready events
func Ready(session *discordgo.Session, event *discordgo.Ready) {
	if strings.Contains(config.Get().Prefix, " ") {
		logger.Panic("Space in prefix is not allowed. Please remove space.")
	}

	logger.Info(fmt.Sprintf("Prefix set to '%s'", config.Get().Prefix))

	if err := session.UpdateGameStatus(0, fmt.Sprintf("%shelp", config.Get().Prefix)); err != nil {
		logger.Warn(fmt.Sprintf("Error updating status: %s", err.Error()))
	}
	logger.Info(fmt.Sprintf("Logged in as user %s#%s(%s)", session.State.User.Username, session.State.User.Discriminator, session.State.User.ID))

	logger.Info("Initializing Audio Engine...")

	initializer.InitAudioEngine(event)
}
