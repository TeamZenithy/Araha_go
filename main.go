package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/bwmarrin/discordgo"
)

var token string
var lavalink *audioengine.Lavalink
var player *audioengine.Player

func init() {
	token = "bottoken"
	// flag.StringVar(&token, "token", "", "token=unprefixed token")
}

func main() {
	flag.Parse()

	if token == "" {
		panic("no token specified!")
	}
	token = "Bot " + token

	dg, err := discordgo.New(token)
	if err != nil {
		panic(err)
	}
	dg.SyncEvents = false

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(voiceServerUpdate)

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("discordgo ready!")
	s.UpdateStatus(0, "gavalink")

	lavalink = audioengine.NewLavalink("1", event.User.ID)

	err := lavalink.AddNodes(audioengine.NodeConfig{
		REST:      "http://localhost:2333",
		WebSocket: "ws://localhost:2333",
		Password:  "youshallnotpass",
	})
	if err != nil {
		log.Println(err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "?!join" {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			log.Println("fail find channel")
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			log.Println("fail find guild")
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				log.Println("trying to connect to channel")
				err = s.ChannelVoiceJoinManual(c.GuildID, vs.ChannelID, false, false)
				if err != nil {
					log.Println(err)
				} else {
					log.Println("channel voice join succeeded")
				}
			}
		}
	} else if strings.HasPrefix(m.Content, "?!play") {
		query := m.Content[7:]
		node, errBestNode := lavalink.BestNode()
		if errBestNode != nil {
			log.Println(errBestNode)
			return
		}
		log.Println(query)
		tracks, errLoadTracks := node.LoadTracks("ytsearch", query)
		if errLoadTracks != nil {
			log.Println(errLoadTracks)
			return
		}
		if tracks.Type != audioengine.TrackLoaded {
			log.Println("weird tracks type", tracks.Type)
		}
		if tracks.Type == audioengine.NoMatches {
			fmt.Println("NO Result")
			return
		}
		track := tracks.Tracks[0].Data
		errPlay := player.Play(track)
		if errPlay != nil {
			log.Println(errPlay)
			return
		}
	} else if m.Content == "?!stop" {
		err := player.Stop()
		if err != nil {
			log.Println(err)
		}
	} else if m.Content == "?!pause" {
		err := player.Pause(!player.Paused())
		if err != nil {
			log.Println(err)
		}
	} else if strings.HasPrefix(m.Content, "?!volume") {
		query := m.Content[9:]
		vol, err := strconv.Atoi(query)
		if err != nil {
			log.Println(err)
			return
		}
		err = player.Volume(vol)
		if err != nil {
			log.Println(err)
		}
	}
}

func voiceServerUpdate(s *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	log.Println("received Vioce Server Update event.")
	vsu := audioengine.VoiceServerUpdate{
		Endpoint: event.Endpoint,
		GuildID:  event.GuildID,
		Token:    event.Token,
	}

	if p, err := lavalink.GetPlayer(event.GuildID); err == nil {
		err = p.Forward(s.State.SessionID, vsu)
		if err != nil {
			log.Println(err)
		}
		return
	}

	node, err := lavalink.BestNode()
	if err != nil {
		log.Println(err)
		return
	}

	handler := new(audioengine.DummyEventHandler)
	player, err = node.CreatePlayer(event.GuildID, s.State.SessionID, vsu, handler)
	if err != nil {
		log.Println(err)
		return
	}
}
