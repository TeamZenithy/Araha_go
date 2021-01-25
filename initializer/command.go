package initializer

import (
	"github.com/TeamZenithy/Araha/commands/general"
	"github.com/TeamZenithy/Araha/commands/music"
	"github.com/TeamZenithy/Araha/commands/setting"
	"github.com/TeamZenithy/Araha/handler"
)

//InitCommands inits command structures
func InitCommands() {
	handler.InitCommands()
	handler.AddCommand(general.Ping())
	handler.AddCommand(general.Help())
	handler.AddCommand(setting.SetLang())
	handler.AddCommand(music.Play())
	handler.AddCommand(music.Pause())
	handler.AddCommand(music.Resume())
	handler.AddCommand(music.NowPlaying())
	handler.AddCommand(music.Queue())
	handler.AddCommand(music.Stop())
	handler.AddCommand(music.Skip())
	handler.AddCommand(music.Seek())
	handler.AddCommand(music.Shuffle())
	handler.AddCommand(music.Volume())
}
