package events

import (
	"log"
	"strconv"

	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/model"
)

type EventHandler interface {
	OnTrackEnd(player *audioengine.Player, track string, reason string) error
	OnTrackException(player *audioengine.Player, track string, reason string) error
	OnTrackStuck(player *audioengine.Player, track string, threshold int) error
}

type EventHandlerManager struct {
	handler []EventHandler
}

func NewEventHandlerManager() *EventHandlerManager {
	log.Println("Added")
	return &EventHandlerManager{
		handler: make([]EventHandler, 0),
	}
}

func (h *EventHandlerManager) OnTrackEnd(player *audioengine.Player, track string, reason string) (err error) {
	if ms, ok := model.Music[player.GuildID()]; ok {
		ms.Queue = ms.Queue[1:]
		if len(ms.Queue) != 0 {
			ms.SongEnd <- "next"
		} else {
			ms.SongEnd <- "end"
			ms.Queue = make([]model.Song, 0)
		}
	}
	return
}

func (h *EventHandlerManager) OnTrackException(player *audioengine.Player, track string, reason string) (err error) {
	log.Printf("Track exception for %s: %s", player.GuildID(), reason)
	if ms, ok := model.Music[player.GuildID()]; ok {
		ms.SongEnd <- reason
	}
	return
}

func (h *EventHandlerManager) OnTrackStuck(player *audioengine.Player, track string, threshold int) (err error) {
	if ms, ok := model.Music[player.GuildID()]; ok {
		ms.SongEnd <- "resume:" + strconv.Itoa(threshold)
	}
	return
}

func (h *EventHandlerManager) AddHandler(handler EventHandler) {
	h.handler = append(h.handler, handler)
}
