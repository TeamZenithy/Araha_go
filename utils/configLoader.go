package utils

import (
	"github.com/TeamZenithy/Araha/structure"
	"log"

	"github.com/BurntSushi/toml"
)

//LoadConfig get rawConfig string and returns error, prefix, token
func loadConfig(rawConfig string) (errLoadFailed error, prefix string, token string) {
	var config structure.Config
	_, err := toml.Decode(rawConfig, &config)
	if err != nil {
		log.Fatal(err)
		return err, "", ""
	}
	return nil, config.Prefix, config.Token
}

//GetToken returns token
func GetToken(rawConfig string) (err error, token string) {
	err, _, token = loadConfig(rawConfig)
	return err, token
}

//GetPrefix returns token
func GetPrefix(rawConfig string) (err error, prefix string) {
	err, prefix, _ = loadConfig(rawConfig)
	return err, prefix
}
