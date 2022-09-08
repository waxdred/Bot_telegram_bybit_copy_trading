package main

import (
	"bot/bybits/bybit"
	"bot/bybits/get"
	"bot/bybits/listen"
	"bot/bybits/post"
	"bot/bybits/telegram"
	"bot/env"
	"fmt"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func run(updates tgbotapi.UpdatesChannel, order *bybit.Bot, trade *bybit.Trades, api env.Env, botapi *tgbotapi.BotAPI) {
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := update.Message.Text
			url := fmt.Sprint(
				"https://api.telegram.org/bot",
				api.Api_telegram,
				"/sendMessage?text=",
				msg,
				"&chat_id=@trading_bybit_wax",
			)
			http.Get(url)

		} else if update.ChannelPost != nil {
			msg := update.ChannelPost.Text
			log.Println(msg)
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
				}
				trd := bybit.GetTrade(dataBybite.Currency, trade)
				if trd != nil {
					px := get.GetPrice(dataBybite.Currency, api)
					sl := post.CancelBySl(px, trd)
					if sl != "" {
						lsErr := post.ChangeLs(api, dataBybite.Currency, sl, trd.Type)
						if lsErr != nil {
							log.Println(lsErr)
						} else {
							log.Printf("Cancel Position ok")
						}
					}
					log.Println(cancelErr)
				}
				trade.Delete(dataBybite.Currency)
				order.Delete(dataBybite.Currency)
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
	err := env.LoadEnv(&api)
	if err != nil {
		log.Println(err)
		log.Fatalf("Error cannot Read file .env")
	}
	log.Println(api)
	order.NewBot(&trade, false)
	log.Printf("Get api Ok")
	botapi, err := tgbotapi.NewBotAPI(api.Api_telegram)
	if err != nil {
		log.Panic(err)
	}

	botapi.Debug = order.Debeug

	log.Printf("Authorized on account %s", botapi.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botapi.GetUpdatesChan(u)

	go listen.GetPositionOrder(api, &trade, &order)
	run(updates, &order, &trade, api, botapi)
}
