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
		LavalinkConfig = make([]string, 0)
	}
	Token = config.Token
	Prefix = config.Prefix
	Owners = config.Owners
	LavalinkConfig = []string{config.LavalinkHost, config.LavalinkPort, config.LavalinkPass}
}
