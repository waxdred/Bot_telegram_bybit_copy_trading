package listen

import (
	"bot/bybits/get"
	"bot/bybits/post"
	"bot/bybits/print"
	"bot/bybits/sign"
	"bot/data"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetPosition(api data.BybitApi, symbol string, url_bybite string) (get.Position, error) {
	var position get.Position
	params := map[string]string{
		"api_key":   api.Api,
		"symbol":    symbol,
		"timestamp": print.GetTimestamp(),
	}
	params["sign"] = sign.GetSigned(params, api.Api_secret)
	url := fmt.Sprint(
		url_bybite,
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
		log.Println(err)
		return position, err
	}
	json.Unmarshal(body, &position)
	return position, nil
}

func BuyTp(api data.BybitApi, trade *data.Trades, symbol string, order *data.Bot, url_bybite string) error {
	price := get.GetPrice(symbol, url_bybite)
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
		err = post.ChangeLs(api, symbol, trade.GetTp2(symbol), trade.GetType(symbol), url_bybite)
		if err == nil {
			trade.SetSl(symbol, trade.GetTp2(symbol))
			log.Printf("%s: Tp2 😎", symbol)
		}
	} else if lastPrice >= tp1 {
		err = post.ChangeLs(api, symbol, trade.GetTp1(symbol), trade.GetType(symbol), url_bybite)
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

func SellTp(api data.BybitApi, trade *data.Trades, symbol string, order *data.Bot, url_bybite string) error {
	price := get.GetPrice(symbol, url_bybite)
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
		err = post.ChangeLs(api, symbol, trade.GetTp2(symbol), trade.GetType(symbol), url_bybite)
		if err == nil {
			trade.SetSl(symbol, trade.GetTp2(symbol))
			log.Printf("%s: Tp2 😎", symbol)
		}
	} else if lastPrice <= tp1 {
		err = post.ChangeLs(api, symbol, trade.GetTp1(symbol), trade.GetType(symbol), url_bybite)
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

func GetPositionOrder(api *data.Env, order *data.Bot) {
	for ok := true; ok; {
		for _, o := range (*order).Active {
			log.Println(o)
		}
		log.Println()
		for _, ord := range (*order).Active {
			log.Print("len order ok:")
			for _, apis := range api.Api {
				log.Print("api: while")
				pos, err := GetPosition(apis, ord.Symbol, api.Url)
				if err == nil {
					order.CheckPositon(pos)
					if apis.Trade.GetType(ord.Symbol) == "Sell" && ord.Active == true {
						err := SellTp(apis, &apis.Trade, ord.Symbol, order, api.Url)
						if err != nil {
							log.Println(err)
						}
					} else if apis.Trade.GetType(ord.Symbol) == "Buy" && ord.Active == true {
						err := BuyTp(apis, &apis.Trade, ord.Symbol, order, api.Url)
						if err != nil {
							log.Println(err)
						}
					}
				} else {
					log.Println(err)
				}
				if order.Debeug {
					log.Println(print.PrettyPrint(apis.Trade))
					log.Println(print.PrettyPrint(order))
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func UpdateChannel(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		}
	}
}
