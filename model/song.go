package model

import (
	audioengine "github.com/TeamZenithy/Araha/engine/audio"
)

//Music is a map of MusicStruct
var Music map[string]*MusicStruct

// MusicStruct is the struct to use in the Music map
type MusicStruct struct {
	ChannelID string
	Queue     []Song
	// "next": skip song, "end": end playback, "resume:TIME", resumes from time (used when track is stuck)
	SongEnd chan string

	PlayerCreated chan bool
	Player        *audioengine.Player
}

// Song is the struct to use in a queue
type Song struct {
	Requester string
	Track     audioengine.Track
	Skips     float64
}
