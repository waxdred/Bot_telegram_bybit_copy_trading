package main

import (
	"bot/bybits/bot"
	"bot/bybits/get"
	"bot/bybits/listen"
	"bot/bybits/post"
	"bot/bybits/telegram"
	"bot/data"
	"bot/mysql"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func run(updates tgbotapi.UpdatesChannel, order *data.Bot, api data.Env) {
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := update.Message.Text
			bot.BotParseMsg(msg, update.Message.From.UserName, &api, order, update)
			log.Println(msg)
			dataBybite, err := telegram.ParseMsg(msg, order.Debeug)
			if err == nil && dataBybite.Trade {
				price := get.GetPrice(dataBybite.Currency, api.Url)
				if price.RetCode == 0 && price.Result[0].BidPrice != "" {
					for _, apis := range api.Api {
						if apis.Trade.Add(apis, dataBybite, price, api.Url) {
							post.PostIsoled(apis, dataBybite.Currency, &apis.Trade, api.Url, order.Debeug)
							err = post.PostOrder(dataBybite.Currency, apis, &apis.Trade, api.Url, order.Debeug)
							if err != nil {
								log.Println(err)
								apis.Trade.Delete(dataBybite.Currency)
							} else {
								order.AddActive(dataBybite.Currency)
							}
						} else {
							if order.Debeug {
								log.Printf("You trade already this Symbol")
							}
						}
						if order.Debeug {
							apis.Trade.Print()
						}
					}
				} else {
					log.Printf("Symbol not found")
				}
			} else if err == nil && dataBybite.Cancel {
				for _, apis := range api.Api {
					cancelErr := post.CancelOrder(dataBybite.Currency, apis, &apis.Trade, api.Url)
					if cancelErr != nil {
						log.Println(cancelErr)
					}
					trd := data.GetTrade(dataBybite.Currency, &apis.Trade)
					if trd != nil {
						px := get.GetPrice(dataBybite.Currency, api.Url)
						sl := post.CancelBySl(px, trd)
						if sl != "" {
							lsErr := post.ChangeLs(apis, dataBybite.Currency, sl, trd.Type, api.Url)
							if lsErr != nil {
								log.Println(lsErr)
							} else {
								log.Printf("Cancel Position ok")
							}
						}
						log.Println(cancelErr)
					}
					apis.Trade.Delete(dataBybite.Currency)
					order.Delete(dataBybite.Currency)
				}
			} else if order.Debeug {
				log.Printf("Error Parsing")
			}
		}
	}
}

func main() {
	var api data.Env
	var order data.Bot

	// waiting mysql running
	log.Print("waiting mysql....")
	time.Sleep(10 * time.Second)

	// for show debeug set at true
	// get var env in struct
	err := data.LoadEnv(&api)
	if err != nil {
		log.Fatal("Error cannot Read file .env: ", err)
	}

	// get data sql set struct order
	if order.NewBot(&api, false) != nil {
		log.Fatalf("NewBot error: ")
	}
	err = mysql.ConnectionDb(&order, &api)
	if err != nil {
		log.Fatal(err)
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

	go listen.GetPositionOrder(&api, &order)
	run(order.Updates, &order, api)
}
