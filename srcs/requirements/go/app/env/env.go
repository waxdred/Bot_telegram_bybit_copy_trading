package env

import (
	"errors"
	"os"
)

type Env struct {
	Api          string
	Api_secret   string
	Api_telegram string
	Url          string
}

func GetEnv(env *Env) error {
	env.Api = os.Getenv("API")
	if env.Api == "" {
		return errors.New("Api not found")
	}
	env.Api_secret = os.Getenv("API_SECRET")
	if env.Api_secret == "" {
		return errors.New("Api_secret not found")
	}
	env.Api_telegram = os.Getenv("API_TELEGRAM")
	if env.Api_telegram == "" {
		return errors.New("Api_telegram not found")
	}
	env.Url = os.Getenv("URL")
	if env.Url == "" {
		return errors.New("Url not found")
	}
	return nil
}

func LoadEnv(env *Env) error {
	err := GetEnv(env)
	if err != nil {
		return err
	}
	return nil
}
