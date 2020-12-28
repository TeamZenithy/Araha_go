package utils

import (
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/structure"

	"github.com/BurntSushi/toml"
)

//LoadConfig get rawConfig string and returns error, prefix, token
func LoadConfig(rawConfig string) {
	var config structure.Config
	_, err := toml.Decode(rawConfig, &config)
	if err != nil {
		logger.Fatal(err.Error())
		Token = ""
		Prefix = ""
		Owners = make([]string, 0)
		ShardLogChannel = ""
		ShardStatusLogChannel = ""
		LavalinkConfig = make([]string, 0)
		RedisHost = ""
		RedisPort = ""
		RedisPass = ""
	}
	Token = config.Token
	Prefix = config.Prefix
	Owners = config.Owners
	ShardLogChannel = config.ShardLogChannel
	ShardStatusLogChannel = config.ShardStatusLogChannel
	LavalinkConfig = []string{config.LavalinkHost, config.LavalinkPort, config.LavalinkPass}
	RedisHost = config.RedisHost
	RedisPort = config.RedisPort
	RedisPass = config.RedisPass
}
