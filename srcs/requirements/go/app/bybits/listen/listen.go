package listen

import (
	"bot/bybits/bybit"
	"bot/bybits/get"
	"bot/bybits/post"
	"bot/bybits/print"
	"bot/bybits/sign"
	"bot/env"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func GetPosition(api env.Env, trade *bybit.Trades) (get.Position, error) {
	var position get.Position
	params := map[string]string{
		"api_key":   api.Api,
		"symbol":    "HNTUSDT",
		"timestamp": print.GetTimestamp(),
	}
	params["sign"] = sign.GetSigned(params, api.Api_secret)
	url := fmt.Sprint(
		api.Url,
		"/private/linear/position/list?api_key=",
		params["api_key"],
		"&symbol=HNTUSDT",
		"&timestamp=",
		params["timestamp"],
		"&sign=",
		params["sign"],
	)
	log.Println(print.PrettyPrint(url))
	body, err := get.GetRequetJson(url)
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(body, &position)
	return position, nil
}

func BuyTp(api env.Env, trade *bybit.Trades, symbol string, order *bybit.Bot) error {
	price := get.GetPrice(symbol, api)
	var err error

	err = nil
	if price.Result[0].LastPrice <= trade.GetSl(symbol) {
		trade.Delete(symbol)
		order.Delete(symbol)
		log.Printf("%s: All take-profit targets achieved 😎", symbol)
	} else if price.Result[0].LastPrice >= trade.GetTp2(symbol) {
		err = post.ChangeLs(api, symbol, trade.GetTp2(symbol))
		log.Printf("%s: Tp2 😎", symbol)
	} else if price.Result[0].LastPrice >= trade.GetTp1(symbol) {
		err = post.ChangeLs(api, symbol, trade.GetTp1(symbol))
		log.Printf("%s: Tp1 😎", symbol)
	}
	if err != nil {
		return err
	}
	return nil
}

func SellTp(api env.Env, trade *bybit.Trades, symbol string, order *bybit.Bot) error {
	price := get.GetPrice(symbol, api)
	var err error

	err = nil
	if price.Result[0].LastPrice >= trade.GetSl(symbol) {
		trade.Delete(symbol)
		order.Delete(symbol)
		log.Printf("%s: All take-profit targets achieved 😎", symbol)
	} else if price.Result[0].LastPrice <= trade.GetTp2(symbol) {
		err = post.ChangeLs(api, symbol, trade.GetTp2(symbol))
		log.Printf("%s: Tp2 😎", symbol)
	} else if price.Result[0].LastPrice <= trade.GetTp1(symbol) {
		err = post.ChangeLs(api, symbol, trade.GetTp1(symbol))
		log.Printf("%s: Tp1 😎", symbol)
	}
	if err != nil {
		return err
	}
	return nil
}

func GetPositionOrder(api env.Env, trade *bybit.Trades, order *bybit.Bot) {
	for ok := true; ok; {
		for i := 0; i < len((*order).Active); i++ {
			if trade.GetType((*order).Active[i]) == "Sell" {
				err := SellTp(api, trade, (*order).Active[i], order)
				if err != nil {
					log.Println(err)
				}
			} else if trade.GetType((*order).Active[i]) == "Buy" {
				err := BuyTp(api, trade, (*order).Active[i], order)
				if err != nil {
					log.Println(err)
				}
			}
		}
		if order.Debeug {
			log.Println(print.PrettyPrint(trade))
			log.Println(print.PrettyPrint(order))
		}
		time.Sleep(10 * time.Second)
	}
}
