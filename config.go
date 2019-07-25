package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type YamlConfig struct {
	Access   ConfigAccess   `yaml:"access"`
	Database ConfigDatabase `yaml:"database"`
}

type ConfigAccess struct {
	Users []*User `yaml:"users"`
}

type ConfigDatabase struct {
	DSN        string `yaml:"dsn"`
	Name       string `yaml:"name"`
	Collection string `yaml:"collection"`
}

func (cfg *YamlConfig) initConfig() error {
	configFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return err
	}

	users = cfg.Access.Users

	return nil
}
