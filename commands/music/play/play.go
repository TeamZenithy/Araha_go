package play

import (
	"fmt"
	"strconv"
	"strings"

	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/extensions/embed"
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
			Description:          &handler.Description{ReqPermsission: "SPEAK", Usage: "play [song title or link]\nplay [song name] $soundcloud"},
		},
	)
}

const (
	commandName = "play"
	commandArg  = "query"
)

func run(ctx handler.CommandContext) error {
	e := embed.New(ctx.Session, ctx.Message.ChannelID)
	guild, err := ctx.Session.State.Guild(ctx.Message.GuildID)
	if err != nil {
		return nil
	}

	query := ctx.Arguments[commandArg]
	if query == "" {
		e.SendEmbed(embed.BADREQ, ctx.T("music:BRUrl"))
		return nil
	}

	var userVoiceState discordgo.VoiceState
	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Message.Author.ID {
			userVoiceState = *vs
		}
	}
	if userVoiceState.UserID == "" {
		e.SendEmbed(embed.BADREQ, ctx.T("music:NotInVoiceChannel"))
		return nil
	}

	node, err := utils.Lavalink.BestNode()
	var tracks *audioengine.Tracks
	var errLoadTracks error
	if strings.HasPrefix(query, "http") && strings.Contains(query, "://") {
		query = strings.ReplaceAll(query, " ", "%20")
		if strings.Contains(query, "twitch.tv") {
			e.SendEmbed(embed.BADREQ, ctx.T("music:BRTwitch"))
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
			e.SendEmbed(embed.BADREQ, ctx.T("music:BRVoiceChannel"))
			return nil
		}
	}

	if err != nil {
		e.SendEmbed(embed.ERR_BOT, ctx.T("music:ErrFind")+"\n"+err.Error())
		return nil
	}

	if errLoadTracks != nil {
		e.SendEmbed(embed.ERR_BOT, ctx.T("music:ErrQuery"+"\n"+err.Error()))
		return nil
	}

	if tracks.Type != audioengine.TrackLoaded {
		if tracks.Type == audioengine.LoadFailed {
			e.SendEmbed(embed.ERR_BOT, ctx.T("music:ErrLoad"))
			return nil
		} else if tracks.Type == audioengine.NoMatches {
			e.SendEmbed(embed.BADREQ, ctx.T("music:BRSearch"))
			return nil
		}
	}

	if tracks.Type == audioengine.PlaylistLoaded {
		for pos := range tracks.Tracks {
			errLoadTracks := queueSong(ctx, tracks.Tracks[pos], ms, len(tracks.Tracks))
			if errLoadTracks != nil {
				logger.Warn(errLoadTracks.Error())
				e.SendEmbed(embed.ERR_BOT, ctx.T("error:ErrAddQueue")+errLoadTracks.Error())
				return nil
			}
		}
		e.SendEmbed(embed.INFO, ctx.T("music:AddedPlaylistQueue", fmt.Sprint(len(tracks.Tracks))))
	} else {
		track := tracks.Tracks[0]
		err = queueSong(ctx, track, ms, 1)
		e.SendEmbed(embed.INFO, ctx.T("music:AddedQueue", track.Info.Title))
		if err != nil {
			_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, ctx.T("error:ErrAddQueue")+err.Error())
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
	e := embed.New(ctx.Session, ctx.Message.ChannelID)
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
	e.SendEmbed(embed.INFO, fmt.Sprintf(":musical_note: Now playing: **%s** ", song.Track.Info.Title), embed.AddProfileFooter(user, ctx.T))

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
