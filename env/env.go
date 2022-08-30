package env

import (
	"os"

	"github.com/joho/godotenv"
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
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}
	GetEnv(env)
	return nil
}
