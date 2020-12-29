package utils

import (
	audioengine "github.com/TeamZenithy/Araha/engine/audio"
	"github.com/TeamZenithy/Araha/lang"
	"github.com/go-redis/redis/v8"
)

//Prefix for global
var Prefix string = ""

//Token for global
var Token string = ""

//Owners for global
var Owners []string = []string{""}

//ShardLogChannel for global
var ShardLogChannel = ""

//ShardStatusLogChannel for global
var ShardStatusLogChannel = ""

//LavalinkConfig for global
var LavalinkConfig []string = []string{""}

//RedisHost for global
var RedisHost = ""

//RedisPort for global
var RedisPort = ""

//RedisPass for global
var RedisPass = ""

//Lavalink is not type
var Lavalink *audioengine.Lavalink

//Player is not type
var Player *audioengine.Player

var RDB *redis.Client

var TR *lang.Tr

const (
	//QUERY_TYPE_YOUTUBE is for ytsearch
	QUERY_TYPE_YOUTUBE = "ytsearch:"
	//QUERY_TYPE_SOUNDCLOUD is for ytsearch
	QUERY_TYPE_SOUNDCLOUD = "scsearch:"
	//QUERY_TYPE_URL is for url search
	QUERY_TYPE_URL = ""
)

//Command Category
const (
	CATEGORY_GENERAL = 0
	CATEGORY_MUSIC   = 1
	CATEGORY_DEV     = 2
	CATEGORY_LOCALE  = 3
)
