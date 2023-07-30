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

func run(updates tgbotapi.UpdatesChannel, order *data.Bot, api *data.Env, trade *data.Trades) {
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := update.Message.Text
			bot.BotParseMsg(msg, update.Message.From.UserName, api, order, update)
			dataBybite, err := telegram.ParseMsg(msg, order.Debeug)
			if err == nil && dataBybite.Trade {
				price := get.GetPrice(dataBybite.Currency, api.Url)
				if price.RetCode == 0 && price.Result[0].BidPrice != "" {
					for _, apis := range api.Api {
						if trade.Add(apis, dataBybite, price, api.Url) == true {
							post.PostIsoled(apis, dataBybite.Currency, trade, api.Url, order.Debeug)
							err = post.PostOrder(dataBybite.Currency, apis, trade, api.Url, order.Debeug)
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
					}
				} else {
					log.Printf("Symbol not found")
				}
			} else if err == nil && dataBybite.Cancel {
				for _, apis := range api.Api {
					cancelErr := post.CancelOrder(dataBybite.Currency, apis, trade, api.Url)
					if cancelErr != nil {
						log.Println(cancelErr)
					} else {
						trade.Delete(dataBybite.Currency)
						order.Delete(dataBybite.Currency)
					}
					trd := data.GetTrade(dataBybite.Currency, trade)
					if trd != nil {
						px := get.GetPrice(dataBybite.Currency, api.Url)
						sl := post.CancelBySl(px, trd)
						if sl != "" {
							lsErr := post.ChangeLs(apis, dataBybite.Currency, sl, trd.Type, api.Url)
							if lsErr != nil {
								log.Println(lsErr)
							} else {
								log.Printf("Change sl for cancel position ok")
							}
						}
						log.Println(cancelErr)
					}
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
	var trade data.Trades

	// waiting mysql running
	log.Print("waiting mysql....")
	time.Sleep(6 * time.Second)

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

	go listen.GetPositionOrder(&api, &order, &trade)
	run(order.Updates, &order, &api, &trade)
}
