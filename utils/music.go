package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/TeamZenithy/Araha/model"
	"github.com/bwmarrin/discordgo"
)

//GetUsersInVoice returns number of users exclude of bot self.
func GetUsersInVoice(guild *discordgo.Guild) int {
	usersInVoice := 0
	for _, vs := range guild.VoiceStates {
		if model.Music[guild.ID].ChannelID == vs.ChannelID {
			usersInVoice++
		}
	}
	return usersInVoice - 1
}

// IsInVoiceWithMusic checks if a user is in the same voice channel as where music is playing
func IsInVoiceWithMusic(guild *discordgo.Guild, userID string) bool {
	for _, vs := range guild.VoiceStates {
		if model.Music[guild.ID].ChannelID == vs.ChannelID && userID == vs.UserID {
			return true
		}
	}
	return false
}

// LeaveAndDestroy leaves the voice channel and destroys the player and queue
func LeaveAndDestroy(s *discordgo.Session, guildID string) string {
	delete(model.Music, guildID)
	if err := Player.Destroy(); err != nil {
		return fmt.Sprintf("Error destroying player for %s: %s", guildID, err)
	}
	if err := s.ChannelVoiceJoinManual(guildID, "", false, true); err != nil {
		return fmt.Sprintf("Couldn't leave voice channel for %s: %s", guildID, err)
	}
	return ""
}

func TimeFormat(duration time.Duration) string {
	return strings.Replace(strings.Replace(strings.Replace(time.Until(time.Now().Add(duration*time.Millisecond)).Round(time.Second).String(), "h", "h ", -1), "m", "m ", -1), "s", "s ", -1)
}
