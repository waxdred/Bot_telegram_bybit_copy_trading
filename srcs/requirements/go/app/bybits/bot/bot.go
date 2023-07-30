package bot

import (
	"bot/data"
	"bot/mysql"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendMsg(msg string, user string, api data.Env) error {
	url := fmt.Sprint(
		"https://api.telegram.org/bot",
		api.Api_telegram,
		"/sendMessage")
	params := map[string]string{
		"text":    msg,
		"chat_id": fmt.Sprint("@", user),
	}
	json_data, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
	} else {
		boby, reqErr := http.Post(url, "application/json", bytes.NewBuffer(json_data))
		if reqErr != nil {
			log.Println(reqErr)
		}
		log.Println(boby.Body)
	}
	return nil
}

func BotParseMsg(
	msg string,
	user string,
	api *data.Env,
	order *data.Bot,
	update tgbotapi.Update,
) error {
	sending := true
	for _, adm := range api.Admin {
		if msg == "/help" && user == adm {
			sender := fmt.Sprint("/add      api:api_secret\n",
				"/delete  api\n",
				"/addAdmin login\n",
				"/deleteAdmin login")
			msgs := tgbotapi.NewMessage(update.Message.Chat.ID, sender)
			if sending == true {
				sending = false
				order.Botapi.Send(msgs)
			}
		} else if strings.Index(msg, " ") == 4 && strings.Compare(msg, "/add ") == 1 && user == adm {
			if strings.Index(msg, ":") < 0 {
				msgs := tgbotapi.NewMessage(update.Message.Chat.ID, "Error try again \n/add api:api_secret")
				if sending == true {
					sending = false
					order.Botapi.Send(msgs)
				}
			} else {
				var msgs tgbotapi.MessageConfig
				api_bybit := msg[strings.Index(msg, " ")+1 : strings.Index(msg, ":")]
				api_secret := msg[strings.Index(msg, ":")+1:]
				api.AddApi(api_bybit, api_secret)
				// add api to database
				if mysql.CheckApi("db.api", order.Db, api_bybit) == true {
					mysql.InsertApi(api_bybit, api_secret, "api", order.Db)
					msgs = tgbotapi.NewMessage(update.Message.Chat.ID, "Api add")
				} else {
					msgs = tgbotapi.NewMessage(update.Message.Chat.ID, "Api are already add")
				}
				if sending == true {
					sending = false
					order.Botapi.Send(msgs)
				}
			}
		} else if strings.Index(msg, " ") == 7 && strings.Compare(msg, "/delete ") == 1 && user == adm {
			msg = msg[strings.Index(msg, " ")+1:]
			sender := api.Delette(msg)
			err := mysql.DbDelete("api", msg, order.Db)
			if err != nil {
				msgs := tgbotapi.NewMessage(update.Message.Chat.ID, sender)
				if sending == true {
					sending = false
					order.Botapi.Send(msgs)
				}
			} else {
				msgs := tgbotapi.NewMessage(update.Message.Chat.ID, sender)
				if sending == true {
					sending = false
					order.Botapi.Send(msgs)
				}
			}
		} else if strings.Index(msg, " ") == 9 && strings.Compare(msg, "/addAdmin ") == 1 && user == adm {
			admin := msg[strings.Index(msg, " ")+1:]
			api.AddAdmin(admin)
			// add api to database
			if mysql.CheckAdmin("db.admin", order.Db, admin) == true {
				mysql.InsertAdmin(admin, "admin", order.Db)
				msgs := tgbotapi.NewMessage(update.Message.Chat.ID, "Admin add")
				if sending == true {
					sending = false
					order.Botapi.Send(msgs)
				}
			} else {
				msgs := tgbotapi.NewMessage(update.Message.Chat.ID, "Api are already add")
				if sending == true {
					sending = false
					order.Botapi.Send(msgs)
				}
			}
		} else if strings.Index(msg, " ") == 12 && strings.Compare(msg, "/deleteAdmin ") == 1 && user == adm {
			msg = msg[strings.Index(msg, " ")+1:]
			sender := api.DeletteAdmin(msg)
			err := mysql.DbDeleteAdmin("admin", msg, order.Db)
			if err != nil {
				msgs := tgbotapi.NewMessage(update.Message.Chat.ID, sender)
				if sending == true {
					sending = false
					order.Botapi.Send(msgs)
				}
			} else {
				msgs := tgbotapi.NewMessage(update.Message.Chat.ID, sender)
				if sending == true {
					sending = false
					order.Botapi.Send(msgs)
				}
			}
		} else if sending == true {
			SendMsg(msg, api.IdCHannel, *api)
			sending = false
		}
	}
	return nil
}
