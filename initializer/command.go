package initializer

import (
	"github.com/TeamZenithy/Araha/commands/general/help"
	"github.com/TeamZenithy/Araha/commands/general/ping"
	"github.com/TeamZenithy/Araha/commands/music/join"
	"github.com/TeamZenithy/Araha/commands/music/play"
	"github.com/TeamZenithy/Araha/handler"
)

//InitCommands inits command structures
func InitCommands() {
	// initializer command map
	handler.InitCommands()
	// register commands
	ping.Initialize()
	help.Initialize()
	join.Initialize()
	play.Initialize()
}
