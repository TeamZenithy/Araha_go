package music

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/extensions/embed"
	"github.com/TeamZenithy/Araha/handler"
	m "github.com/TeamZenithy/Araha/middlewares"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

func Play() *handler.Cmd {
	return &handler.Cmd{
		Run:         play,
		Name:        "play",
		Category:    "music",
		Aliases:     []string{"p"},
		Args:        []string{"+url", "?platform"},
		Middlewares: []handler.HandlerFunc{m.UserVoiceState()},
		Usage:       "(URL | Name of Song) --[soundcloud]",
	}
}

func play(c *handler.Context) {
	userVoiceState := c.Get("userVoiceState").(*discordgo.VoiceState)
	urlArg := c.Arg("url")

	node, err := utils.Lavalink.BestNode()
	if err != nil {
		c.Embed.SendEmbed(embed.ERR_BOT, c.T("music:ErrFind")+"\n"+err.Error())
		return
	}
	var tracks *audioengine.Tracks
	var errLoadTracks error
	if strings.HasPrefix(urlArg, "http") && strings.Contains(urlArg, "://") {
		if strings.Contains(urlArg, "twitch.tv") {
			c.Embed.SendEmbed(embed.BADREQ, c.T("music:BRTwitch"))
			return
		}
		urlArg = url.QueryEscape(urlArg)
		tracks, errLoadTracks = node.LoadTracks(utils.QUERY_TYPE_URL, urlArg)
	} else if c.Arg("platform") == "soundcloud" {
		tracks, errLoadTracks = node.LoadTracks(utils.QUERY_TYPE_SOUNDCLOUD, urlArg)
	} else {
		tracks, errLoadTracks = node.LoadTracks(utils.QUERY_TYPE_YOUTUBE, urlArg)
	}

	if errLoadTracks != nil {
		c.Embed.SendEmbed(embed.ERR_BOT, c.T("music:ErrQuery"+"\n"+err.Error()))
		return
	}

	var ms *model.MusicStruct
	inSameVoice, inVoice := false, false
	if temp, ok := model.Music[c.Msg.GuildID]; ok {
		ms = temp
		for _, vs := range c.Msg.Guild().VoiceStates {
			if ms.ChannelID == vs.ChannelID {
				if vs.UserID == c.Msg.Author.ID {
					inSameVoice = true
					break
				}
				inVoice = true
			}
		}
	}

	if !inSameVoice {
		if !inVoice {
			err := c.Session.ChannelVoiceJoinManual(c.Msg.GuildID, userVoiceState.ChannelID, false, true)
			if err != nil {
				return
			}

			model.Music[c.Msg.GuildID] = &model.MusicStruct{
				ChannelID:     userVoiceState.ChannelID,
				Queue:         make([]model.Song, 0),
				SongEnd:       make(chan string),
				PlayerCreated: make(chan bool),
			}
			ms = model.Music[c.Msg.GuildID]
		} else {
			c.Embed.SendEmbed(embed.BADREQ, c.T("music:BRVoiceChannel"))
			return
		}
	}

	if tracks.Type != audioengine.TrackLoaded {
		if tracks.Type == audioengine.LoadFailed {
			c.Embed.SendEmbed(embed.ERR_BOT, c.T("music:ErrLoad"))
			return
		} else if tracks.Type == audioengine.NoMatches {
			c.Embed.SendEmbed(embed.BADREQ, c.T("music:BRSearch"))
			return
		}
	}

	if tracks.Type == audioengine.PlaylistLoaded {
		for pos := range tracks.Tracks {
			errLoadTracks := queueSong(c, tracks.Tracks[pos], ms, len(tracks.Tracks))
			if errLoadTracks != nil {
				c.Embed.SendEmbed(embed.ERR_BOT, c.T("error:ErrAddQueue")+errLoadTracks.Error())
				return
			}
		}
		c.Embed.SendEmbed(embed.INFO, c.T("music:AddedPlaylistQueue", fmt.Sprint(len(tracks.Tracks))))
	} else {
		track := tracks.Tracks[0]
		err = queueSong(c, track, ms, 1)
		c.Embed.SendEmbed(embed.INFO, c.T("music:AddedQueue", track.Info.Title))
		if err != nil {
			_, err = c.Session.ChannelMessageSend(c.Msg.ChannelID, c.T("error:ErrAddQueue")+err.Error())
			return
		}
	}
	return
}

func queueSong(ctx *handler.Context, track audioengine.Track, ms *model.MusicStruct, length int) (err error) {
	justJoined := false
	if len(ms.Queue) == 0 {
		justJoined = true
	}

	ms.Queue = append(ms.Queue, model.Song{
		Requester: ctx.Msg.Author.ID,
		Track:     track,
	})

	if justJoined {
		<-ms.PlayerCreated
		go playSong(ctx, ms.Queue[0], ms, -1)
		if err != nil {
			_, err = ctx.Session.ChannelMessageSend(ctx.Msg.ChannelID, "An error occured during song playback. Please try again.\n"+err.Error())
			return
		}
	}
	return
}

func playSong(ctx *handler.Context, song model.Song, ms *model.MusicStruct, startTime int) (err error) {
	e := embed.New(ctx.Session, ctx.Msg.ChannelID)
	err = utils.Player.Play(song.Track.Data)
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Msg.ChannelID, "Error playing *"+song.Track.Info.Title+"*. Skipping to next song.\n"+err.Error())
		ms.Queue = ms.Queue[1:]
		if len(ms.Queue) != 0 {
			playSong(ctx, ms.Queue[0], ms, -1)
		}
		return
	}
	user, _ := ctx.Session.User(song.Requester)
	e.SendEmbed(embed.INFO, fmt.Sprintf(":musical_note: Now playing: **%s** ", song.Track.Info.Title), embed.AddProfileFooter(user, ctx.T))

	end := <-ms.SongEnd
	if end == "next" {
		err = playSong(ctx, ms.Queue[0], ms, -1)
	} else if end == "end" {
		utils.LeaveAndDestroy(ctx.Session, ctx.Msg.GuildID)
	} else if strings.HasPrefix(end, "resume:") {
		end = end[7:]
		time, err := strconv.Atoi(end)
		if err != nil {
			return err
		}
		err = playSong(ctx, song, ms, time)
	} else {
		_, err = ctx.Session.ChannelMessageSend(ctx.Msg.ChannelID, "재생중 문제가 발생했습니다.")
	}

	return
}
