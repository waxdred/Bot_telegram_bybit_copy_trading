package main

import (
	"bot/bybits/bot"
	"bot/bybits/bybit"
	"bot/bybits/get"
	"bot/bybits/listen"
	"bot/bybits/post"
	"bot/bybits/telegram"
	"bot/env"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func run(updates tgbotapi.UpdatesChannel, order *bybit.Bot, trade *bybit.Trades, api env.Env) {
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := update.Message.Text
			bot.BotParseMsg(msg, update.Message.From.UserName, &api, order, update)
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

	// waiting mysql running
	log.Print("waiting mysql....")
	time.Sleep(10 * time.Second)

	// for show debeug set at true
	// get var env in struct
	err := env.LoadEnv(&api)
	if err != nil {
		log.Fatal("Error cannot Read file .env: ", err)
	}

	// get data sql set struct order
	if order.NewBot(&trade, &api, false) != nil {
		log.Fatalf("NewBot error: ")
	}
	defer order.Db.Close()

	// print api find
	api.ListApi()
	log.Printf("Get api Ok")

	// connection bot telegram
	order.Botapi, err = tgbotapi.NewBotAPI(api.Api_telegram)
	if err != nil {
		log.Panic(err)
	}

	order.Botapi.Debug = order.Debeug

	log.Printf("Authorized on account %s", order.Botapi.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	order.Updates = order.Botapi.GetUpdatesChan(u)

	go listen.GetPositionOrder(api, &trade, &order)
	run(order.Updates, &order, &trade, api)
}
