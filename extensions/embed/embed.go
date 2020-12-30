package embed

import (
	"time"

	"github.com/TeamZenithy/Araha/lang"
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
	T             func(string)
}

const (
	FOOTER = iota + 1
	TITLE
	LINK
	THUMBNAIL
	FIELDS
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
			footer := d.Data.(*discordgo.MessageEmbedFooter)
			embed.Footer = footer
			embed.Timestamp = time.Now().Format(time.RFC3339)
		case TITLE:
			text := d.Data.(string)
			embed.Title = text
		case LINK:
			text := d.Data.(string)
			embed.URL = text
		case THUMBNAIL:
			text := d.Data.(string)
			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: text, Width: 1280, Height: 720}
		case FIELDS:
			fields := d.Data.([]*discordgo.MessageEmbedField)
			embed.Fields = fields
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

func AddTitle(message string) *EmbedMore {
	return &EmbedMore{EmbedMoreType: TITLE, Data: message}
}

func AddLink(link string) *EmbedMore {
	return &EmbedMore{EmbedMoreType: LINK, Data: link}
}

func AddFooter(msg string) *EmbedMore {
	return &EmbedMore{EmbedMoreType: FOOTER, Data: &discordgo.MessageEmbedFooter{Text: msg}}
}

func AddProfileFooter(user *discordgo.User, T lang.HFType) *EmbedMore {
	return &EmbedMore{EmbedMoreType: FOOTER, Data: &discordgo.MessageEmbedFooter{Text: T("music:ReqBy", user.Username, user.Discriminator), IconURL: user.AvatarURL("")}}
}

func AddThumbnail(url string) *EmbedMore {
	return &EmbedMore{EmbedMoreType: THUMBNAIL, Data: url}
}

func AddFields(fields []*discordgo.MessageEmbedField) *EmbedMore {
	return &EmbedMore{EmbedMoreType: FIELDS, Data: fields}
}
