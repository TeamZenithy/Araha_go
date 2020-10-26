package play

import (
	"fmt"
	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/utils"
	"io/ioutil"
	"log"
)

func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Names:                []string{commandName},
			RequiredArgumentType: []string{commandArg},
			Usage:                map[string]string{"필요한 권한":"**``음성 채널 발언권``**", "설명":"``요청된 이름의 노래 또는 링크를 검색해서 음원을 재생합니다.``", "사용법": "```css\n?!ping 노래 이름 또는 링크```"},
		},
	)
}

const (
	commandName = "play"
	commandArg = "없음"

	QUERY_TYPE_YOUTUBE= "ytsearch"
	QUERY_TYPE_SOUNDCLOUD = "scsearch"
)

func run(ctx handler.CommandContext) error {
	rawConfig, errFindConfigFile := ioutil.ReadFile("config.toml") // just pass the file name
	if errFindConfigFile != nil {
		log.Fatalln("Error while load config file: " + errFindConfigFile.Error())
		return nil
	}
	errLoadConfigData, prefix := utils.GetPrefix(string(rawConfig))
	if errLoadConfigData != nil {
		log.Fatalln("Error while load config data: " + errLoadConfigData.Error())
	}
	query := ctx.Message.Content[len(prefix) + len(commandName) + 1:]
	node, errBestNode := utils.Lavalink.BestNode()
	if errBestNode != nil {
		log.Println(errBestNode)
		return nil
	}
	ctx.Message.Reply("🔎 " + query + "을(를) 찾는중...")
	tracks, errLoadTracks := node.LoadTracks("ytsearch", query)
	if errLoadTracks != nil {
		log.Println(errLoadTracks)
		return nil
	}
	if tracks.Type != audioengine.TrackLoaded && tracks.Type != audioengine.SearchResult {
		log.Println("weird tracks type", tracks.Type)
	}
	if tracks.Type == audioengine.NoMatches {
		fmt.Println("NO Result")
		return nil
	}
	track := tracks.Tracks[0].Data
	errPlay := utils.Player.Play(track)
	ctx.Message.Reply("Track Info: ```json\n"+fmt.Sprintf("%b", tracks.Tracks[0].Info)+"```")
	if errPlay != nil {
		log.Println(errPlay)
		return nil
	}
	return nil
}