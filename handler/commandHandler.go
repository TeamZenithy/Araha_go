package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/TeamZenithy/Araha/logger"

	"github.com/TeamZenithy/Araha/extensions/embed"

	"github.com/TeamZenithy/Araha/lang"
	"github.com/TeamZenithy/Araha/utils"

	"github.com/TeamZenithy/Araha/db"

	"github.com/TeamZenithy/Araha/extensions/objects"

	"github.com/bwmarrin/discordgo"
)

type Argument map[string]string
type HandlerFunc func(c *Context) bool

const (
	keywordArgumentPrefix = "--"
	stringSeparator       = " "
)

var (
	//Commands is a map of string
	Commands map[string]*Cmd
	//Aliases is a map of string
	Aliases map[string]string
)

//CommandContext is ctx
type CommandContext struct {
	Session   *discordgo.Session
	Message   *objects.ExtendedMessage
	Arguments map[string]string
	arguments Argument
	T         lang.HFType
	Locale    string
}

type Context struct {
	Session *discordgo.Session
	Msg     *objects.ExtendedMessage
	Embed   *embed.Embed
	args    Argument
	T       lang.HFType
	Locale  string
	values  map[string]interface{}
}

type Cmd struct {
	Run         func(c *Context)
	Middlewares []HandlerFunc
	Name        string
	Category    string
	Aliases     []string
	Args        []string
	Usage       string
}

//Command includes Run(function), Name(list of srings), Required args type(list of strings), and Usage(map of string and string)
type Command struct {
	Run                  func(ctx CommandContext) error
	Name                 string
	Aliases              []string
	RequiredArgumentType []string
	Category             int
	Usage                map[string]interface{}
	Description          *Description
}

type Description struct {
	ReqPermsission string
	Usage          string
}

//InitCommands Initialize the Commands map
func InitCommands() {
	Commands = make(map[string]*Cmd)
	Aliases = make(map[string]string)
}

func (c *Context) Set(key string, value interface{}) {
	c.values[key] = value
}

func (c *Context) Get(key string) interface{} {
	return c.values[key]
}

func (c *Context) Arg(key string) string {
	return c.args[key]
}

//AddCommand Adds a command to the Commands map
func AddCommand(cmd *Cmd) {
	Commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		Aliases[alias] = cmd.Name
	}
	logger.Info("Command Added: " + cmd.Name)
}

func FindCommand(name string) *Cmd {
	if command, ok := Commands[name]; ok {
		return command
	} else if command, ok := Aliases[name]; ok {
		return Commands[command]
	}
	return nil
}

//HandleCreatedMessage Handle a message creation event
func HandleCreatedMessage(s *discordgo.Session, m *discordgo.MessageCreate, prefix string) {
	if m.Author.Bot {
		return
	}
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	list := strings.Split(m.Content, " ")
	cmdName := string([]rune(list[0])[len(prefix):])
	var args []string
	if len(list) > 1 {
		args = list[1:]
	}

	var err error = nil

	command := FindCommand(cmdName)

	if command == nil {
		return
	}

	userLocale := ""
	l, err := db.FindUserLocale(m.Author.ID)
	if err != nil {
		fmt.Println("Error while parsing locale")
		userLocale = "en"
	} else if l == "" {
		l, err := db.FindGuildLocale(m.GuildID)
		if err != nil {
			fmt.Println("Error while parsing locale")
			userLocale = "en"
		}
		userLocale = l
	} else {
		userLocale = l
	}
	T := utils.TR.GetHandlerFunc(userLocale, "en")
	parsed, err := ParseArgument(command.Args, args)
	e := embed.New(s, m.ChannelID)
	if err != nil {
		e.SendEmbed(embed.BADREQ, T("error:ErrSyntex", prefix, command.Name))
		return
	}

	context := &Context{
		Session: s,
		Msg:     objects.ExtendMessage(m.Message, s),
		args:    *parsed,
		Embed:   e,
		T:       T,
		Locale:  userLocale,
		values:  make(map[string]interface{}),
	}

	for _, d := range command.Middlewares {
		if result := d(context); !result {
			return
		}
	}

	command.Run(context)
}

func ParseArgument(arg, content []string) (*Argument, error) {
	parsed := Argument{}

	startIndex := 0

	argLen := len(arg)
	if argLen < 1 {
		return &parsed, nil
	}
	argLast := arg[argLen-1]
	if strings.HasPrefix(argLast, "?") {
		if argLen < 2 {
			if len(content) < 1 {
				return &parsed, nil
			}
			parsed[strings.TrimPrefix(argLast, "?")] = content[0]
			return &parsed, nil
		}
		contentLen := len(content)
		lastItem := content[contentLen-1]
		if strings.HasPrefix(lastItem, "--") {
			argLast = strings.TrimPrefix(argLast, "?")
			parsed[argLast] = strings.ReplaceAll(lastItem, "--", "")
			content = content[0 : contentLen-1]
		}
		arg = arg[0 : argLen-1]
	}

	if len(content) < len(arg) {
		return &Argument{}, errors.New("Error Parse")
	}

	for _, d := range arg {
		if strings.HasPrefix(d, "+") {
			parsed[strings.TrimLeft(d, "+")] = strings.Join(content[startIndex:], " ")
		} else {
			parsed[d] = content[startIndex]
			startIndex++
		}
	}
	return &parsed, nil
}
