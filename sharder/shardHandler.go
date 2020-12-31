package sharder

import (
	"github.com/bwmarrin/discordgo"
)

func (m *Manager) OnDiscordConnected(s *discordgo.Session, evt *discordgo.Connect) {
	m.handleEvent(EventConnected, s.ShardID+1, "")
}

func (m *Manager) OnDiscordDisconnected(s *discordgo.Session, evt *discordgo.Disconnect) {
	m.handleEvent(EventDisconnected, s.ShardID+1, "")
}

func (m *Manager) OnDiscordReady(s *discordgo.Session, evt *discordgo.Ready) {
	m.handleEvent(EventReady, s.ShardID+1, "")
}

func (m *Manager) OnDiscordResumed(s *discordgo.Session, evt *discordgo.Resumed) {
	m.handleEvent(EventResumed, s.ShardID+1, "")
}
