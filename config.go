package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type YamlConfig struct {
	Access      ConfigAccess      `yaml:"access" validate:"required"`
	Database    ConfigDatabase    `yaml:"database" validate:"required"`
	Application ConfigApplication `yaml:"application" validate:"required"`
}

type ConfigAccess struct {
	Users []*User `yaml:"users" validate:"required"`
}

type ConfigDatabase struct {
	DSN        string        `yaml:"dsn" validate:"required,mongodsn"`
	Name       string        `yaml:"name" validate:"required,alphanum"`
	Collection string        `yaml:"collection" validate:"required,alphanum"`
	Timeout    time.Duration `yaml:"timeout" validate:"required,min=1"`
}

type ConfigApplication struct {
	Templates []*Template `yaml:"templates" validate:"required"`
}

type Template struct {
	Name string `yaml:"name" validate:"required,alphanum"`
	Path string `yaml:"path" validate:"required,file"`
}

func (cfg *YamlConfig) initConfig() error {
	configFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(configFile, &cfg); err != nil {
		return err
	}

	if err := validate.Struct(cfg); err != nil {
		return err
	}

	users = cfg.Access.Users

	templates.Delims("#(", ")#")
	_, err = templates.ParseFiles(cfg.getTemplatesPaths()...)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *YamlConfig) getTemplatesPaths() []string {
	result := make([]string, 0)
	for _, t := range cfg.Application.Templates {
		result = append(result, t.Path)
	}
	return result
}
