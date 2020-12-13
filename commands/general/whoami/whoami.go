package whoami

import (
	"fmt"
	"io/ioutil"

	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/utils"

	"github.com/TeamZenithy/Araha/handler"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"wai"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_GENERAL,
			Usage:                map[string]string{"필요한 권한": "**``없음``**", "설명": "``자신이 누구인지 확인합니다.``", "사용법": "```css\n?!whoami```"},
		},
	)
}

const (
	commandName = "whoami"
	commandArg  = "없음"
)

func run(ctx handler.CommandContext) error {
	rawConfig, errFindConfigFile := ioutil.ReadFile("config.toml") // just pass the file name
	if errFindConfigFile != nil {
		logger.Fatal(fmt.Sprintf("Error while load config file: %s", errFindConfigFile.Error()))
		return nil
	}
	owners, errLoadConfigData := utils.GetOwners(string(rawConfig))
	if errLoadConfigData != nil {
		logger.Fatal(fmt.Sprintf("Error while load config data: %s", errLoadConfigData.Error()))
		return nil
	}
	if utils.Contains(owners, ctx.Message.Author.ID) {
		_, err := ctx.Message.Reply("You are owner")
		return err
	}
	_, err := ctx.Message.Reply("You are not owner")

	return err
}
