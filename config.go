package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
)

type Config struct {
	GithubToken string `json:"github-token"`
}

func ReadConfigFromFile() (*Config, error) {
	user, _ := user.Current()
	filePath := user.HomeDir + "/.feature-up.config"

	if _, err := os.Stat(filePath); err == nil {
		var c Config

		raw, jsonErr := ioutil.ReadFile(filePath)
		if jsonErr != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		json.Unmarshal(raw, &c)
		return &c, nil
	} else {
		return nil, errors.New("Could not find config file")
	}
}

func (c *Config) Save() {
	user, _ := user.Current()
	filePath := user.HomeDir + "/.feature-up.config"

	bytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile(filePath, bytes, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
