package listen

import (
	"bybit/bybit/bybit"
	"bybit/bybit/get"
	"bybit/bybit/print"
	"bybit/bybit/sign"
	"bybit/env"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func GetPosition(api env.Env, trade *bybit.Trades) (get.Position, error) {
	var position get.Position
	params := map[string]string{
		"api_key":   api.Api,
		"timestamp": print.GetTimestamp(),
	}
	params["sign"] = sign.GetSigned(params, api.Api_secret)
	url := fmt.Sprint(
		api.Url,
		"/private/linear/position/list?api_key=",
		params["api_key"],
		"&timestamp=",
		params["timestamp"],
		"&sign=",
		params["sign"],
	)
	body, err := get.GetRequetJson(url)
	if err != nil {
		log.Panic(err)
	}
	jsonErr := json.Unmarshal(body, &position)
	if jsonErr != nil {
		log.Panic(err)
	}
	return position, nil
}

func GetPositionOrder(api env.Env, trade *bybit.Trades) {
	for ok := true; ok; {
		time.Sleep(10 * time.Second)
		_, err := GetPosition(api, trade)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Print Position")
	}
}
