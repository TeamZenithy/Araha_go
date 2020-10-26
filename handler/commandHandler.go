package handler

import (
	"github.com/TeamZenithy/Araha/extensions/objects"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ztrue/tracerr"
)

const (
	keywordArgumentPrefix = "--"
	stringSeparator       = " "
)

var (
	Commands map[string]Command
)

type CommandContext struct {
	Session   *discordgo.Session
	Message   *objects.ExtendedMessage
	Arguments map[string]string
}

type Command struct {
	Run                  func(ctx CommandContext) error
	Names                []string
	RequiredArgumentType []string
	Usage                map[string]string
}

// Initialize the Commands map
func InitCommands() {
	Commands = make(map[string]Command)
}

// Adds a command to the Commands map
func AddCommand(command Command) {
	for _, name := range command.Names {
		Commands[name] = command
	}
}

// Handle a message creation event
func HandleCreatedMessage(session *discordgo.Session, message *discordgo.MessageCreate, prefix string) {
	if message.Author.ID == session.State.User.ID || message.Author.Bot {
		return
	}
	if !strings.HasPrefix(message.Content, prefix) {
		return
	}

	var endIndex = strings.Index(message.Content, stringSeparator)

	var commandName string

	if endIndex == -1 {
		commandName = message.Content[len(prefix):]
	} else {
		commandName = message.Content[len(prefix):strings.Index(message.Content, stringSeparator)]
	}

	var command, exists = Commands[commandName]

	if !exists {
		return
	}

	var context = CommandContext{
		Session: session,
		Message: objects.ExtendMessage(message.Message, session),
		Arguments: parseArguments(
			message.Content,
			command.RequiredArgumentType,
			command.Usage),
	}

	var err = command.Run(context)
	if err != nil {
		tracerr.PrintSourceColor(err)
		_, err = session.ChannelMessageSend(
			message.ChannelID,
			fmt.Sprint("An Error occurred while executing the command!\n", err))
		if err != nil {
			tracerr.PrintSourceColor(err)
			_, err = session.ChannelMessageSend(
				message.ChannelID,
				"An Error occurred while executing the command and sending the error message!")
			if err != nil {
				tracerr.PrintSourceColor(err)
			}
		}
	}
}

// Parse command arguments
func parseArguments(
	content string,
	expectedPositionalArguments []string,
	keywordArgumentAliases map[string]string) map[string]string {

	var separated = strings.Split(content, stringSeparator)

	// do not process the command name and prefix
	separated = separated[1:]

	var returnArguments = make(map[string]string)

	// if len(separated) == 0 {
	//	return returnArguments
	// }

	var currentPosition = 0

	for len(separated) > 0 {
		// remove first element from slice
		var currentItem = separated[0]
		separated = separated[1:]

		var currentArgumentValue []string

		if strings.HasPrefix(currentItem, keywordArgumentPrefix) {
			for len(separated) > 0 && !strings.HasPrefix(separated[0], keywordArgumentPrefix) {
				_ = append(currentArgumentValue, separated[0])
				separated = separated[1:]
				currentPosition++
			}
		} else {
			if currentPosition >= len(expectedPositionalArguments) {
				var _, exists = returnArguments[expectedPositionalArguments[len(expectedPositionalArguments)-1]]
				if exists {
					// The length checks should prevent the value from being nil
					//goland:noinspection GoNilness
					returnArguments[expectedPositionalArguments[len(expectedPositionalArguments)-1]] += stringSeparator + currentItem
				} else {
					//goland:noinspection GoNilness
					returnArguments[expectedPositionalArguments[len(expectedPositionalArguments)-1]] = currentItem
				}
			} else {
				//goland:noinspection GoNilness
				returnArguments[expectedPositionalArguments[currentPosition]] = currentItem
			}
		}

		currentPosition++
	}

	//goland:noinspection GoNilness
	for key, value := range returnArguments {
		key = strings.ToLower(key)
		var _, exists = keywordArgumentAliases[key]
		if exists {
			returnArguments[keywordArgumentAliases[key]] = value
		}
	}

	return returnArguments
}
