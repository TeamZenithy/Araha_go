package utils

import (
	"log"

	"github.com/TeamZenithy/Araha/structure"

	"github.com/BurntSushi/toml"
)

//LoadConfig get rawConfig string and returns error, prefix, token
func loadConfig(rawConfig string) (prefix string, token string, lavalinkHost string, lavalinkPort string, lavalinkPass string, errLoadFailed error) {
	var config structure.Config
	_, err := toml.Decode(rawConfig, &config)
	if err != nil {
		log.Fatal(err)
		return "", "", "", "", "", err
	}
	return config.Prefix, config.Token, config.LavalinkHost, config.LavalinkPort, config.LavalinkPass, nil
}

//GetToken returns token
func GetToken(rawConfig string) (token string, err error) {
	_, token, _, _, _, err = loadConfig(rawConfig)
	return token, err
}

//GetPrefix returns prefix
func GetPrefix(rawConfig string) (prefix string, err error) {
	prefix, _, _, _, _, err = loadConfig(rawConfig)
	return prefix, err
}

//GetLavalinkConfig returns lavalink configs
func GetLavalinkConfig(rawConfig string) (host, port, pass string, err error) {
	_, _, host, port, pass, err = loadConfig(rawConfig)
	return host, port, pass, err
}
