package objects

import (
	"errors"
	"fmt"

	"github.com/TeamZenithy/Araha/extensions"
	"github.com/bwmarrin/discordgo"
)

type ExtendedMessage struct {
	*discordgo.Message
	session *discordgo.Session
}

func ExtendMessage(message *discordgo.Message, session *discordgo.Session) *ExtendedMessage {
	return &ExtendedMessage{
		Message: message,
		session: session,
	}
}

// short form for message.session.Guild(message.GuildID)
func (message *ExtendedMessage) Guild() *ExtendedGuild {
	guild, err := message.session.State.Guild(message.GuildID)
	if err != nil {
		return nil
	}
	return ExtendGuild(guild, message.session)
}

func (message *ExtendedMessage) Member() *ExtendedMember {
	member, err := message.session.GuildMember(message.GuildID, message.ID)
	if err != nil {
		return nil
	}
	return ExtendMember(member, message.session)
}

// short form for message.session.ChannelMessageSend(message.ChannelID, content)
func (message *ExtendedMessage) Reply(content string) (*discordgo.Message, error) {
	return message.session.ChannelMessageSend(message.ChannelID, content)
}

// short form for message.session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{...})
func (message *ExtendedMessage) ComplexReply(send *discordgo.MessageSend) (*discordgo.Message, error) {
	return message.session.ChannelMessageSendComplex(message.ChannelID, send)
}

// short form for message.session.Channel(message.ChannelID)
func (message *ExtendedMessage) Channel() (*discordgo.Channel, error) {
	return message.session.Channel(message.ChannelID)
}

func (message *ExtendedMessage) AuthorMember() (*discordgo.Member, error) {
	messageGuild := message.Guild()
	if messageGuild == nil {
		return nil, errors.New("")
	}
	guildMembers, err := messageGuild.GetMembers()
	if err != nil {
		return nil, err
	}

	for _, member := range guildMembers {
		if member.User.ID == message.Author.ID {
			return member, nil
		}
	}
	return nil, errors.New(fmt.Sprint(extensions.MemberNotFoundError,
		"member ", message.Author.ID, " not found in guild ", message.GuildID))
}
