package main

import (
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

func addressListNameValidator(fl validator.FieldLevel) bool {

	if ok, _ := regexp.MatchString(`^[A-Za-z0-9-]+$`, fl.Field().String()); !ok {
		return false
	}

	return true
}

func commentValidator(fl validator.FieldLevel) bool {

	if ok, _ := regexp.MatchString(`^[A-Za-zА-Яа-я0-9\s,.:-]+$`, fl.Field().String()); !ok {
		return false
	}

	return true
}

func accessKeyValidator(fl validator.FieldLevel) bool {

	if ok, _ := regexp.MatchString(`^[A-Z0-9]{24}$`, fl.Field().String()); !ok {
		return false
	}

	return true
}

func secretKeyValidator(fl validator.FieldLevel) bool {

	if ok, _ := regexp.MatchString(`^[a-f0-9]{64}$`, fl.Field().String()); !ok {
		return false
	}

	return true
}

func mongoDSNValidator(fl validator.FieldLevel) bool {

	if ok, _ := regexp.MatchString(`^mongodb://([a-zA-Z0-9][a-zA-Z0-9.-]+[a-zA-Z0-9]|(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])):([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$`, fl.Field().String()); !ok {
		return false
	}

	return true
}

func registerValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("addresslistname", addressListNameValidator); err != nil {
		return err
	}

	if err := v.RegisterValidation("comment", commentValidator); err != nil {
		return err
	}

	if err := v.RegisterValidation("accesskey", accessKeyValidator); err != nil {
		return err
	}

	if err := v.RegisterValidation("secretkey", secretKeyValidator); err != nil {
		return err
	}

	if err := v.RegisterValidation("mongodsn", mongoDSNValidator); err != nil {
		return err
	}

	return nil
}
