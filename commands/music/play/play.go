package play

import (
	"fmt"
	"strconv"
	"strings"

	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/handler"
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
			Usage:                map[string]string{"필요한 권한": "**``음성 채널 발언권``**", "설명": "``요청된 이름의 노래 또는 링크를 검색해서 음원을 재생합니다.``", "사용법": fmt.Sprintf("```css\n%splay 노래 이름 또는 링크```", utils.Prefix)},
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

	query := strings.Replace(strings.Join(strings.Split(ctx.Message.Content, " ")[1:], " "), " ", "%20", -1)
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

	node, err := utils.Lavalink.BestNode()
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error finding music node. Please try again.\n"+err.Error())
		return nil
	}

	var tracks *audioengine.Tracks
	var errLoadTracks error
	if !strings.HasPrefix(query, "http") && !strings.Contains(query, "://") {
		tracks, errLoadTracks = node.LoadTracks(utils.QUERY_TYPE_YOUTUBE, query)
	} else {
		tracks, errLoadTracks = node.LoadTracks(utils.QUERY_TYPE_URL, query)
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

	track := tracks.Tracks[0]
	err = queueSong(ctx, track, ms)
	if err != nil {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error adding song to queue. Please try again.\n"+err.Error())
		return nil
	}

	return nil
}

func queueSong(ctx handler.CommandContext, track audioengine.Track, ms *model.MusicStruct) (err error) {
	justJoined := false
	if len(ms.Queue) == 0 {
		justJoined = true
	}

	ms.Queue = append(ms.Queue, model.Song{
		Requester: ctx.Message.Author.Username,
		Track:     track,
	})

	if justJoined {
		<-ms.PlayerCreated
		err = playSong(ctx, ms.Queue[0], ms, -1)
		if err != nil {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "An error occured during song playback. Please try again.\n"+err.Error())
			return
		}
	} else {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Queued for playback: **"+track.Info.Title+"**")
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
	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf(":musical_note: Now playing: **%s** *(requested by %s)*", song.Track.Info.Title, song.Requester))

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
