package utils

import (
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/structure"

	"github.com/BurntSushi/toml"
)

//LoadConfig get rawConfig string and returns error, prefix, token
func loadConfig(rawConfig string) (prefix string, token string, owners []string, lavalinkHost string, lavalinkPort string, lavalinkPass string, errLoadFailed error) {
	var config structure.Config
	_, err := toml.Decode(rawConfig, &config)
	if err != nil {
		logger.Fatal(err.Error())
		return "", "", make([]string, 0), "", "", "", err
	}
	return config.Prefix, config.Token, config.Owners, config.LavalinkHost, config.LavalinkPort, config.LavalinkPass, nil
}

//GetToken returns token
func GetToken(rawConfig string) (token string, err error) {
	_, token, _, _, _, _, err = loadConfig(rawConfig)
	return token, err
}

//GetPrefix returns prefix
func GetPrefix(rawConfig string) (prefix string, err error) {
	prefix, _, _, _, _, _, err = loadConfig(rawConfig)
	return prefix, err
}

//GetLavalinkConfig returns lavalink configs
func GetLavalinkConfig(rawConfig string) (host, port, pass string, err error) {
	_, _, _, host, port, pass, err = loadConfig(rawConfig)
	return host, port, pass, err
}

//GetOwners returns prefix
func GetOwners(rawConfig string) (owners []string, err error) {
	_, _, owners, _, _, _, err = loadConfig(rawConfig)
	return owners, err
}
