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
	"strconv"
	"time"
)

func GetPosition(api env.Env, trade *bybit.Trades, symbol string) (get.Position, error) {
	var position get.Position
	params := map[string]string{
		"api_key":   api.Api,
		"symbol":    symbol,
		"timestamp": print.GetTimestamp(),
	}
	params["sign"] = sign.GetSigned(params, api.Api_secret)
	url := fmt.Sprint(
		api.Url,
		"/private/linear/position/list?api_key=",
		params["api_key"],
		"&symbol=", symbol,
		"&timestamp=",
		params["timestamp"],
		"&sign=",
		params["sign"],
	)
	body, err := get.GetRequetJson(url)
	if err != nil {
		log.Panic(err)
	}
	json.Unmarshal(body, &position)
	// log.Println(print.PrettyPrint(position))
	return position, nil
}

func BuyTp(api env.Env, trade *bybit.Trades, symbol string, order *bybit.Bot) error {
	price := get.GetPrice(symbol, api)
	lastPrice, _ := strconv.ParseFloat(price.Result[0].LastPrice, 64)
	sl, _ := strconv.ParseFloat(trade.GetSl(symbol), 64)
	tp1, _ := strconv.ParseFloat(trade.GetTp1(symbol), 64)
	tp2, _ := strconv.ParseFloat(trade.GetTp2(symbol), 64)
	tp3, _ := strconv.ParseFloat(trade.GetTp3(symbol), 64)
	var err error

	err = nil
	if lastPrice <= sl {
		trade.Delete(symbol)
		order.Delete(symbol)
		log.Printf("%s: Sl touch: ", symbol)
	} else if lastPrice >= tp3 {
		log.Printf("%s: All take-profit targets achieved 😎: ", symbol)
		trade.Delete(symbol)
		order.Delete(symbol)
	} else if lastPrice >= tp2 {
		err = post.ChangeLs(api, symbol, trade.GetTp2(symbol), trade.GetType(symbol))
		if err == nil {
			trade.SetSl(symbol, trade.GetTp2(symbol))
			log.Printf("%s: Tp2 😎", symbol)
		}
	} else if lastPrice >= tp1 {
		err = post.ChangeLs(api, symbol, trade.GetTp1(symbol), trade.GetType(symbol))
		if err == nil {
			trade.SetSl(symbol, trade.GetTp1(symbol))
			log.Printf("%s: Tp1 😎", symbol)
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func SellTp(api env.Env, trade *bybit.Trades, symbol string, order *bybit.Bot) error {
	price := get.GetPrice(symbol, api)
	lastPrice, _ := strconv.ParseFloat(price.Result[0].LastPrice, 64)
	sl, _ := strconv.ParseFloat(trade.GetSl(symbol), 64)
	tp1, _ := strconv.ParseFloat(trade.GetTp1(symbol), 64)
	tp2, _ := strconv.ParseFloat(trade.GetTp2(symbol), 64)
	tp3, _ := strconv.ParseFloat(trade.GetTp3(symbol), 64)
	var err error

	err = nil
	if lastPrice >= sl {
		trade.Delete(symbol)
		order.Delete(symbol)
		log.Printf("%s: Sl touch: ", symbol)
	} else if lastPrice <= tp3 {
		trade.Delete(symbol)
		order.Delete(symbol)
		log.Printf("%s: All take-profit targets achieved 😎: ", symbol)
	} else if lastPrice <= tp2 {
		err = post.ChangeLs(api, symbol, trade.GetTp2(symbol), trade.GetType(symbol))
		if err == nil {
			trade.SetSl(symbol, trade.GetTp2(symbol))
			log.Printf("%s: Tp2 😎", symbol)
		}
	} else if lastPrice <= tp1 {
		err = post.ChangeLs(api, symbol, trade.GetTp1(symbol), trade.GetType(symbol))
		if err == nil {
			trade.SetSl(symbol, trade.GetTp1(symbol))
			log.Printf("%s: Tp1 😎", symbol)
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func GetPositionOrder(api env.Env, trade *bybit.Trades, order *bybit.Bot) {
	for ok := true; ok; {
		for i := 0; i < len((*order).Active); i++ {
			pos, _ := GetPosition(api, trade, (*order).Active[i].Symbol)
			order.CheckPositon(pos)
			if trade.GetType((*order).Active[i].Symbol) == "Sell" && (*order).Active[i].Active == true {
				err := SellTp(api, trade, (*order).Active[i].Symbol, order)
				if err != nil {
					log.Println(err)
				}
			} else if trade.GetType((*order).Active[i].Symbol) == "Buy" && (*order).Active[i].Active == true{
				err := BuyTp(api, trade, (*order).Active[i].Symbol, order)
				if err != nil {
					log.Println(err)
				}
			}
		}
		if order.Debeug {
			log.Println(print.PrettyPrint(trade))
			log.Println(print.PrettyPrint(order))
		}
		time.Sleep(5 * time.Second)
	}
}
