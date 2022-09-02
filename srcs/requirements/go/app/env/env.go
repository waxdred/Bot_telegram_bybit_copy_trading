package env

import (
	"os"
)

type Env struct {
	Api          string `json:"api"`
	Api_secret   string `json:"api_secret"`
	Api_telegram string `json:"api_telegram"`
	Url          string `json:"url"`
}

func GetEnv(env *Env) {
	env.Api = os.Getenv("API")
	env.Api_secret = os.Getenv("API_SECRET")
	env.Api_telegram = os.Getenv("API_TELEGRAM")
	env.Url = os.Getenv("URL")
}

func LoadEnv(env *Env) error {
	GetEnv(env)
	return nil
}
