package config

import (
	"github.com/BurntSushi/toml"
	"github.com/TeamZenithy/Araha/logger"
)

var config ConfigStruct

//Config for discord bot.
type ConfigStruct struct {
	Prefix                string
	Token                 string
	Owners                []string
	ShardLogChannel       string
	ShardStatusLogChannel string
	LavalinkHost          string
	LavalinkPort          string
	LavalinkPass          string
	RedisHost             string
	RedisPort             string
	RedisPass             string
	Release               bool
}

func LoadConfig() {
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		logger.Fatal(err.Error())
	}
}

func Get() *ConfigStruct {
	return &config
}
