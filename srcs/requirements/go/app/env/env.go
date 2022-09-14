package env

import (
	"bot/bybits/print"
	"errors"
	"log"
	"os"
)

type BybitApi struct {
	Api        string
	Api_secret string
}

type Env struct {
	Api          []BybitApi
	Api_telegram string
	Url          string
}

func (t *Env) AddApi(api string, api_secret string) {
	elem := BybitApi{
		Api:        api,
		Api_secret: api_secret,
	}
	(*t).Api = append((*t).Api, elem)
}

func (t *Env) Delette(api string) string {
	ret := false
	ls := (*t).Api
	var tmp []BybitApi

	for i := 0; i < len(ls); i++ {
		if ls[i].Api != api {
			tmp = append(tmp, ls[i])
		} else {
			ret = true
		}
	}
	(*t).Api = tmp
	if ret == false {
		return "Api not found cannot be deletted"
	}
	return "Api deletted"
}

func (t Env) ListApi() {
	for i := 0; i < len(t.Api); i++ {
		log.Println(print.PrettyPrint(t.Api[i]))
	}
}

func GetEnv(env *Env) error {
	api := os.Getenv("API")
	if api == "" {
		return errors.New("Api not found")
	}
	api_secret := os.Getenv("API_SECRET")
	if api_secret == "" {
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
	env.AddApi(api, api_secret)
	return nil
}

func LoadEnv(env *Env) error {
	err := GetEnv(env)
	if err != nil {
		return err
	}
	return nil
}
