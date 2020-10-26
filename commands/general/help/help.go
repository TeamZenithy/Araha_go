package help

import (
	"github.com/TeamZenithy/Araha/handler"
	"fmt"
	"strings"
)

func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Names:                []string{commandName},
			RequiredArgumentType: []string{commandArg},
			Usage:                map[string]string{"필요한 권한":"**``없음``**", "설명":"``모든 명령어의 도움말을 표시합니다.``", "사용법": "```css\n?!help 명령어```"},
		},
	)
}

const (
	commandName = "help"
	commandArg = "명령어"
)

func run(ctx handler.CommandContext) error {
	if _, exists := ctx.Arguments[commandArg]; exists {
		// show specific information about a command
		var command, exists_ = handler.Commands[strings.ToLower(ctx.Arguments[commandArg])]
		if exists_ {
			var formattedCommandNames []string

			for _, value := range command.Names {
				formattedCommandNames = append(
					formattedCommandNames,
					fmt.Sprint(">>> 명령어: `", value, "`"))
			}

			var formattedRequiredArgumentType []string

			for _, value := range command.RequiredArgumentType {
				formattedRequiredArgumentType = append(
					formattedRequiredArgumentType,
					fmt.Sprint("`", value, "`"))
			}

			var formattedUsage []string

			for key, value := range command.Usage {
				formattedUsage = append(
					formattedUsage,
					fmt.Sprint("", key, ": ", value, "\n"))
			}

			var _, err = ctx.Message.Reply(
				fmt.Sprint(
					strings.Join(formattedCommandNames, ", "),
					"\n필요로 하는 인자: ",
					strings.Join(formattedRequiredArgumentType, ", "),
					"\n",
					strings.Join(formattedUsage, "")))
			return err
		} else {
			// command doesn't exist
			var _, err = ctx.Message.Reply("해당 명령어를 찾을 수 없습니다!")
			return err
		}
	} else {
		// list commands
		var outputStr = ">>> 명령어 목록:\n"
		for commandName := range handler.Commands {
			outputStr += fmt.Sprint("`", commandName, "`, ")
		}
		outputStr = outputStr[:len(outputStr)-len(", ")]
		var _, err = ctx.Message.Reply(outputStr)
		return err
	}
}