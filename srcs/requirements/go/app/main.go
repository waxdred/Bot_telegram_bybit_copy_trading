package main

import (
	"bybit/bybit/bybit"
	"bybit/bybit/get"
	"bybit/bybit/listen"
	"bybit/bybit/post"
	"bybit/bybit/telegram"
	"bybit/env"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func run(updates tgbotapi.UpdatesChannel, order *bybit.Bot, trade *bybit.Trades, api env.Env) {
	for update := range updates {
		if update.ChannelPost != nil {
			msg := update.ChannelPost.Text
			dataBybite, err := telegram.ParseMsg(msg, order.Debeug)
			if err == nil && dataBybite.Trade {
				price := get.GetPrice(dataBybite.Currency, api)
				if price.RetCode == 0 && price.Result[0].BidPrice != "" {
					if trade.Add(api, dataBybite, price) {
						post.PostIsoled(api, dataBybite.Currency, trade, order.Debeug)
						err = post.PostOrder(dataBybite.Currency, api, trade, order.Debeug)
						if err != nil {
							log.Println(err)
							trade.Delete(dataBybite.Currency)
						} else {
							order.AddActive(dataBybite.Currency)
						}
					} else {
						if order.Debeug {
							log.Printf("You trade already this Symbol")
						}
					}
					if order.Debeug {
						trade.Print()
					}
				} else {
					log.Printf("Symbol not found")
				}
			} else if err == nil && dataBybite.Cancel {
				cancelErr := post.CancelOrder(dataBybite.Currency, api, trade)
				if cancelErr != nil {
					log.Println(cancelErr)
				} else if order.Debeug {
					trade.Delete(dataBybite.Currency)
					order.Delete(dataBybite.Currency)
				}
			} else if order.Debeug {
				log.Printf("Error Parsing")
			}
		}
	}
}

func main() {
	var api env.Env
	var order bybit.Bot
	var trade bybit.Trades

	// for show debeug set at true
	order.NewBot(&trade, false)
	err := env.LoadEnv(&api)
	if err != nil {
		log.Fatalf("Error cannot Read file .env")
	}
	log.Printf("Get api Ok")
	botapi, err := tgbotapi.NewBotAPI(api.Api_telegram)
	if err != nil {
		log.Panic(err)
	}

	botapi.Debug = false

	log.Printf("Authorized on account %s", botapi.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botapi.GetUpdatesChan(u)

	go listen.GetPositionOrder(api, &trade, &order)
	run(updates, &order, &trade, api)
}
