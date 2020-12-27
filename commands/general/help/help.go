package help

import (
	"fmt"
	"strings"

	"github.com/TeamZenithy/Araha/utils"

	"github.com/TeamZenithy/Araha/handler"
)

//Initialize command
func Initialize() {
	handler.AddCommand(
		handler.Command{
			Run:                  run,
			Name:                 commandName,
			Aliases:              []string{"h", "guide", "manual"},
			RequiredArgumentType: []string{commandArg},
			Category:             utils.CATEGORY_GENERAL,
			Usage:                map[string]string{"Required Permission": "**``none``**", "Description": "``모든 명령어의 도움말을 표시합니다.``", "Usage": fmt.Sprintf("```css\n%shelp [command]```", utils.Prefix)},
		},
	)
}

const (
	commandName = "help"
	commandArg  = "Command"
)

func run(ctx handler.CommandContext) error {
	if _, exists := ctx.Arguments[commandArg]; exists {
		// show specific information about a command
		command, isExists := handler.Commands[strings.ToLower(ctx.Arguments[commandArg])]
		aliasCommand, isExistsAliasCommand := handler.Aliases[strings.ToLower(ctx.Arguments[commandArg])]

		if !isExists && isExistsAliasCommand {
			isExists = isExistsAliasCommand
			command = handler.Commands[aliasCommand]
		}

		if isExists {
			var formattedCommandNames []string

			formattedCommandNames = append(
				formattedCommandNames,
				fmt.Sprint(">>> Command: `", command.Name, "`"))

			var formattedRequiredArgumentType []string

			for _, value := range command.RequiredArgumentType {
				formattedRequiredArgumentType = append(
					formattedRequiredArgumentType,
					fmt.Sprint("`", value, "`"))
			}

			var formattedCommandAliases []string

			if formattedCommandAliases == nil || len(command.Aliases) == 1 && command.Aliases[0] == "" {
				formattedCommandAliases = append(formattedCommandAliases, fmt.Sprint("`none`"))
			} else {
				for _, value := range command.Aliases {
					formattedCommandAliases = append(formattedCommandAliases, fmt.Sprint("`", value, "`"))
				}
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
					"\nAlias: ",
					strings.Join(formattedCommandAliases, ", "),
					"\nRequired Argument(s): ",
					strings.Join(formattedRequiredArgumentType, ", "),
					"\n",
					strings.Join(formattedUsage, "")))
			return err
		}
		var _, err = ctx.Message.Reply("The command was not found!")
		return err
	}
	var outputStr = ">>> List of commands:\n"
	for commandName := range handler.Commands {
		outputStr += fmt.Sprint("`", commandName, "`, ")
	}
	outputStr = outputStr[:len(outputStr)-len(", ")]
	var _, err = ctx.Message.Reply(outputStr)
	return err
}
