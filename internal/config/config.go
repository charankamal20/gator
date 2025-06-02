package config

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() Config {
	data, err := os.ReadFile(getConfigFilePath())
	if err != nil {
		panic(err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func (conf Config) PrintConfig() {
	if conf.CurrentUsername == "" {
		println("No user set")
	} else {
		println("Current user:", conf.CurrentUsername)
	}

	println("DB URL:", conf.DBUrl)
}

func (conf Config) SetUser(username string) {
	conf.CurrentUsername = username
	data, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(getConfigFilePath(), data, 0644)
	if err != nil {
		panic(err)
	}
}


func getConfigFilePath() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic("Could not determine home directory")
	}

	return path.Join(homePath, configFileName)
}
