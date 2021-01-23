package initializer

import (
	"github.com/TeamZenithy/Araha/commands/general/help"
	"github.com/TeamZenithy/Araha/commands/general/ping"
	"github.com/TeamZenithy/Araha/commands/general/whoami"
	"github.com/TeamZenithy/Araha/commands/music/nowplaying"
	"github.com/TeamZenithy/Araha/commands/music/pause"
	"github.com/TeamZenithy/Araha/commands/music/play"
	"github.com/TeamZenithy/Araha/commands/music/queue"
	"github.com/TeamZenithy/Araha/commands/music/resume"
	"github.com/TeamZenithy/Araha/commands/music/seek"
	"github.com/TeamZenithy/Araha/commands/music/shuffle"
	"github.com/TeamZenithy/Araha/commands/music/skip"
	"github.com/TeamZenithy/Araha/commands/music/stop"
	"github.com/TeamZenithy/Araha/commands/music/volume"
	"github.com/TeamZenithy/Araha/commands/setting/locale"
	"github.com/TeamZenithy/Araha/handler"
)

//InitCommands inits command structures
func InitCommands() {
	// initializer command map
	handler.InitCommands()
	// register commands
	ping.Initialize()
	help.Initialize()
	whoami.Initialize()
	play.Initialize()
	skip.Initialize()
	queue.Initialize()
	stop.Initialize()
	seek.Initialize()
	volume.Initialize()
	shuffle.Initialize()
	nowplaying.Initialize()
	pause.Initialize()
	resume.Initialize()
	locale.Initialize()
}
