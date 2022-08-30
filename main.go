package main

import (
	"bybit/bybit/bybit"
	"bybit/bybit/get"
	"bybit/bybit/listen"
	"bybit/bybit/post"
	"bybit/bybit/print"
	"bybit/bybit/telegram"
	"bybit/env"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	var api env.Env
	var trade bybit.Trades
	// var trade map[string]bybit.Trade

	err := env.LoadEnv(&api)
	if err != nil {
		log.Fatalf("Error cannot Read file .env")
	}
	log.Printf("Get api Ok")
	botapi, err := tgbotapi.NewBotAPI(api.Api_telegram)
	if err != nil {
		log.Panic(err)
	}

	botapi.Debug = true

	log.Printf("Authorized on account %s", botapi.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botapi.GetUpdatesChan(u)

	go listen.GetPosition(api, &trade)

	for update := range updates {
		if update.ChannelPost != nil {
			msg := update.ChannelPost.Text
			dataBybite, err := telegram.ParseMsg(msg)
			if err == nil && dataBybite.Trade {
				price := get.GetPrice(dataBybite.Currency, api)
				if price.RetCode == 0 {
					if trade.Add(api, dataBybite, price) {
						err = post.PostOrder(dataBybite.Currency, api, &trade)
						if err != nil {
							log.Println(err)
						}
					} else {
						log.Printf("You trade already this Symbol")
					}
					trade.Print()
				}
			} else if err == nil && dataBybite.Cancel {
				post.CancelOrder(dataBybite.Currency, api, &trade)
				log.Printf("Cancel: %s", dataBybite.Currency)
				log.Println(print.PrettyPrint(trade))
			} else {
				log.Printf("Error Parsing")
			}
		}
	}
}
