package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"

	valid "mikrotik_provisioning/pkg/validator"
)

const (
	configFile = "config.yml"
)

type (
	Config struct {
		Access      *Access      `yaml:"access" validator:"required"`
		DB          *Database    `yaml:"database" validator:"required"`
		Application *Application `yaml:"application" validator:"required"`
	}

	Access struct {
		Users []*User `yaml:"users" validator:"required"`
	}

	User struct {
		AccessKey string `yaml:"access_key" bson:"access_key" validator:"required,access_key"`
		SecretKey string `yaml:"secret_key" bson:"secret_key" validator:"required,secret_key"`
	}

	Database struct {
		DSN         string        `yaml:"dsn" validator:"required,mongo_dsn"`
		Name        string        `yaml:"name" validator:"required,alphanum"`
		Collections []*Collection `yaml:"collections" validator:"required"`
		Timeout     time.Duration `yaml:"timeout" validator:"required,min=1"`
	}

	Collection struct {
		Resource string               `yaml:"resource" validator:"required,alphanum"`
		Name     string               `yaml:"name" validator:"required,alphanum"`
		Indexes  []*CollectionIndexes `yaml:"indexes" validator:"required"`
	}

	CollectionIndexes struct {
		Name   string `yaml:"name" validator:"required,alphanum"`
		Unique bool   `yaml:"unique" validator:"required"`
		Field  string `yaml:"field" validator:"required,alphanum"`
	}

	Application struct{}

	Template struct {
		Name string `yaml:"name" validator:"required,alphanum"`
		Path string `yaml:"path" validator:"required,file"`
	}
)

func ParseConfig() (*Config, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	validator := validator.New()
	if err := valid.RegisterValidators(validator); err != nil {
		return nil, err
	}
	if err := validator.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}
