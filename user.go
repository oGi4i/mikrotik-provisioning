package main

type User struct {
	AccessKey string `yaml:"access_key" bson:"access_key" validate:"required,accesskey"`
	SecretKey string `yaml:"secret_key" bson:"secret_key" validate:"required,secretkey"`
}
