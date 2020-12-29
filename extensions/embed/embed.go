package embed

import (
	"github.com/bwmarrin/discordgo"
)

type EmbedType int
type EmbedMoreType int

type Embed struct {
	ChannelID string
	Session   *discordgo.Session
}

type EmbedMore struct {
	EmbedMoreType EmbedMoreType
	Data          interface{}
}

const (
	FOOTER = iota + 1
	FOOTER_USER
)

const (
	ERR_BOT = iota + 1
	BADREQ
	INFO
	ETC
)

func New(session *discordgo.Session, channelID string) *Embed {
	return &Embed{ChannelID: channelID, Session: session}
}

func (e *Embed) SendEmbed(embedType EmbedType, message string, embedMore ...*EmbedMore) {
	embed := discordgo.MessageEmbed{Description: message}
	for _, d := range embedMore {
		switch d.EmbedMoreType {
		case FOOTER:
			text := d.Data.(string)
			embed.Footer = &discordgo.MessageEmbedFooter{Text: text}
		}
	}
	switch embedType {
	case ERR_BOT:
		embed.Color = 0xffd5cd
	case BADREQ:
		embed.Color = 0xefbbcf
	case INFO:
		embed.Color = 0xc3aed6
	case ETC:
		break
	default:
		return
	}
	e.Session.ChannelMessageSendEmbed(e.ChannelID, &embed)
}

func AddFooter(message string) *EmbedMore {
	return &EmbedMore{EmbedMoreType: FOOTER, Data: message}
}
