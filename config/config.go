package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	valid "mikrotik_provisioning/validate"
	"time"
)

const (
	configFile = "config.yml"
)

var (
	Config = &YamlConfig{}
)

type YamlConfig struct {
	Access      *AccessConfig      `yaml:"access" validate:"required"`
	Database    *DatabaseConfig    `yaml:"database" validate:"required"`
	Application *ApplicationConfig `yaml:"application" validate:"required"`
}

type AccessConfig struct {
	Users []*User `yaml:"users" validate:"required"`
}

type User struct {
	AccessKey string `yaml:"access_key" bson:"access_key" validate:"required,accesskey"`
	SecretKey string `yaml:"secret_key" bson:"secret_key" validate:"required,secretkey"`
}

type DatabaseConfig struct {
	DSN         string                 `yaml:"dsn" validate:"required,mongodsn"`
	Name        string                 `yaml:"name" validate:"required,alphanum"`
	Collections []*CollectionMapConfig `yaml:"collections" validate:"required"`
	Timeout     time.Duration          `yaml:"timeout" validate:"required,min=1"`
}

type CollectionMapConfig struct {
	Resource string                     `yaml:"resource" validate:"required,alphanum"`
	Name     string                     `yaml:"name" validate:"required,alphanum"`
	Indexes  []*CollectionIndexesConfig `yaml:"indexes" validate:"required"`
}

type CollectionIndexesConfig struct {
	Name   string `yaml:"name" validate:"required,alphanum"`
	Unique bool   `yaml:"unique" validate:"required"`
	Field  string `yaml:"field" validate:"required,alphanum"`
}

type ApplicationConfig struct {
}

type Template struct {
	Name string `yaml:"name" validate:"required,alphanum"`
	Path string `yaml:"path" validate:"required,file"`
}

func (cfg *YamlConfig) InitConfig() error {
	configFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(configFile, &cfg); err != nil {
		return err
	}

	if err := valid.Validate.Struct(cfg); err != nil {
		return err
	}

	return nil
}
