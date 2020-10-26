package play

import (
	"fmt"
	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/handler"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Names:                []string{commandName},
			RequiredArgumentType: []string{commandArg},
			Usage:                map[string]string{"필요한 권한": "**``음성 채널 발언권``**", "설명": "``요청된 이름의 노래 또는 링크를 검색해서 음원을 재생합니다.``", "사용법": "```css\n?!ping 노래 이름 또는 링크```"},
		},
	)
}

const (
	commandName = "play"
	commandArg  = "노래 이름 또는 링크"

	QUERY_TYPE_YOUTUBE    = "ytsearch"
	QUERY_TYPE_SOUNDCLOUD = "scsearch"
)

func run(ctx handler.CommandContext) error {
	node, errBestNode := utils.Lavalink.BestNode()
	if errBestNode != nil {
		log.Println(errBestNode)
		return nil
	}
	if ctx.Arguments[commandArg] == "" {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "❌ 검색어 또는 링크를 입력해주세요.")
		return nil
	}
	searchingMsg, _ := ctx.Message.Reply("🔎 " + ctx.Arguments[commandArg] + "을(를) 찾는중...")
	tracks, errLoadTracks := node.LoadTracks(QUERY_TYPE_YOUTUBE, ctx.Arguments[commandArg])
	if errLoadTracks != nil {
		log.Println(errLoadTracks)
		return nil
	}
	if tracks.Type != audioengine.TrackLoaded && tracks.Type != audioengine.SearchResult {
		log.Println("weird tracks type", tracks.Type)
	}
	if tracks.Type == audioengine.NoMatches {
		fmt.Println("NO Result")
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "❌ 일치하는 검색 결과가 없습니다.")
		return nil
	}

	track := tracks.Tracks[0].Data

	errPlay := utils.Player.Play(track)
	if errPlay != nil {
		log.Println(errPlay)
		return nil
	}
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  39423,
		Title:  "✅ 노래가 추가되었습니다.",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "제목",
				Value:  tracks.Tracks[0].Info.Title,
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:   "업로더",
				Value:  tracks.Tracks[0].Info.Author,
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:   "링크",
				Value:  tracks.Tracks[0].Info.URI,
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:   "신청자",
				Value:  "<@" + ctx.Message.Author.ID + ">",
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	ctx.Session.ChannelMessageDelete(searchingMsg.ChannelID, searchingMsg.ID)
	ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)

	return nil
}
