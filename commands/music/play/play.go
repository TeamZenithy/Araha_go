package play

import (
	"fmt"
	"strconv"
	"strings"

	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/model"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"p"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_MUSIC,
			Usage:                map[string]string{"Required Permission": "**``SPEAK``**", "Description": "``Search for the requested song title or link to play the song.\nIf you want search song on soudcloud, then include ''$soundcloud'' in message.``", "Usage": fmt.Sprintf("```css\n%splay [song title or link]```\n```css\n%splay [song name] --soundcloud```", utils.Prefix, utils.Prefix)},
		},
	)
}

const (
	commandName = "play"
	commandArg  = "노래 이름 또는 링크"
)

func run(ctx handler.CommandContext) error {
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil
	}

	query := ctx.Arguments[commandArg]
	if query == "" {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Please provide a url or search query.")
		return nil
	}

	var userVoiceState discordgo.VoiceState
	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			userVoiceState = *vs
		}
	}
	if userVoiceState.UserID == "" {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "You are not in a voice channel.")
		return nil
	}

	node, err := utils.Lavalink.BestNode()
	var tracks *audioengine.Tracks
	var errLoadTracks error
	if strings.HasPrefix(query, "http") && strings.Contains(query, "://") {
		query = strings.ReplaceAll(query, " ", "%20")
		if strings.Contains(query, "twitch.tv") {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "I'm sorry. It's not good to watch Twitch using a bot. Please use the browser instead of me.")
			return nil
		}
		tracks, errLoadTracks = node.LoadTracks(utils.QUERY_TYPE_URL, query)
	} else if strings.Contains(query, "$soundcloud") {
		query = strings.Replace(query, "$soundcloud", "", -1)
		tracks, errLoadTracks = node.LoadTracks(utils.QUERY_TYPE_SOUNDCLOUD, query)
	} else {
		tracks, errLoadTracks = node.LoadTracks(utils.QUERY_TYPE_YOUTUBE, query)
	}

	var ms *model.MusicStruct
	alreadyInSameVoice, alreadyInVoice := false, false
	if temp, ok := model.Music[ctx.Message.GuildID]; ok {
		ms = temp
		for _, vs := range guild.VoiceStates {
			if ms.ChannelID == vs.ChannelID {
				if vs.UserID == ctx.Message.Author.ID {
					alreadyInSameVoice = true
					break
				}
				alreadyInVoice = true
			}
		}
	}

	if !alreadyInSameVoice {
		if !alreadyInVoice {
			err := ctx.Session.ChannelVoiceJoinManual(ctx.Message.GuildID, userVoiceState.ChannelID, false, true)
			if err != nil {
				return err
			}

			model.Music[ctx.Message.GuildID] = &model.MusicStruct{
				ChannelID:     userVoiceState.ChannelID,
				Queue:         make([]model.Song, 0),
				SongEnd:       make(chan string),
				PlayerCreated: make(chan bool),
			}
			ms = model.Music[ctx.Message.GuildID]
		} else {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "I am already in a different voice channel; not switching.")
			return nil
		}
	}

	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error finding music node. Please try again.\n"+err.Error())
		return nil
	}

	if errLoadTracks != nil {
		_, errLoadTracks = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error with query. Please try again or try a different query.\n"+err.Error())
		return nil
	}

	if tracks.Type != audioengine.TrackLoaded {
		if tracks.Type == audioengine.LoadFailed {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Track failed to load. Please try again.")
			return nil
		} else if tracks.Type == audioengine.NoMatches {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "No matches for that query. Please try a different query.")
			return nil
		}
	}

	if tracks.Type == audioengine.PlaylistLoaded {
		for pos := range tracks.Tracks {
			errLoadTracks := queueSong(ctx, tracks.Tracks[pos], ms, len(tracks.Tracks))
			if errLoadTracks != nil {
				logger.Warn(errLoadTracks.Error())
				ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error adding songs to queue. Please try again.\n"+errLoadTracks.Error())
				return nil
			}
		}
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("Queued **%d** songs...", len(tracks.Tracks)))
	} else {
		track := tracks.Tracks[0]
		err = queueSong(ctx, track, ms, 1)
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("Queued for playback: **%s**", track.Info.Title))
		if err != nil {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error adding song to queue. Please try again.\n"+err.Error())
			return nil
		}
	}

	return nil
}

func queueSong(ctx handler.CommandContext, track audioengine.Track, ms *model.MusicStruct, length int) (err error) {
	justJoined := false
	if len(ms.Queue) == 0 {
		justJoined = true
	}

	ms.Queue = append(ms.Queue, model.Song{
		Requester: ctx.Message.Author.ID,
		Track:     track,
	})

	if justJoined {
		<-ms.PlayerCreated
		go playSong(ctx, ms.Queue[0], ms, -1)
		if err != nil {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "An error occured during song playback. Please try again.\n"+err.Error())
			return
		}
	}
	return
}

func playSong(ctx handler.CommandContext, song model.Song, ms *model.MusicStruct, startTime int) (err error) {
	err = utils.Player.Play(song.Track.Data)
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error playing *"+song.Track.Info.Title+"*. Skipping to next song.\n"+err.Error())
		ms.Queue = ms.Queue[1:]
		if len(ms.Queue) != 0 {
			playSong(ctx, ms.Queue[0], ms, -1)
		}
		return
	}
	user, _ := ctx.Session.User(song.Requester)
	username := user.Username
	discriminator := user.Discriminator
	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(":musical_note: Now playing: **%s** *(requested by %s)*", song.Track.Info.Title, fmt.Sprintf("%s#%s", username, discriminator)))

	end := <-ms.SongEnd
	if end == "next" {
		err = playSong(ctx, ms.Queue[0], ms, -1)
	} else if end == "end" {
		utils.LeaveAndDestroy(ctx.Session, ctx.Message.GuildID)
	} else if strings.HasPrefix(end, "resume:") {
		end = end[7:]
		time, err := strconv.Atoi(end)
		if err != nil {
			return err
		}
		err = playSong(ctx, song, ms, time)
	} else {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "재생중 문제가 발생했습니다.")
	}

	return
}
