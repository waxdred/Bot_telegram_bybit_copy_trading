package telegram

import (
	"bybit/bybit/print"
	"errors"
	"log"
	"strings"
)

type Data struct {
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Entry    string `json:"entry"`
	Tp1      string `json:"tp_1"`
	Tp2      string `json:"tp_2"`
	Tp3      string `json:"tp_3"`
	Sl       string `json:"sl"`
	Level    string `json:"level"`
	SetUp    string `json:"set_up"`
	Order    string `json:"order"`
	Cancel   bool
	Trade    bool
}

func ParseMsg(msg string, debug bool) (Data, error) {
	var data Data

	data.Cancel = false
	data.Trade = false
	if strings.Index(msg, "Cancelled") > 0 {
		pos := strings.Index(msg, "#")
		if pos > 0 {
			data.Cancel = true
			data.Currency = msg[pos+1:]
			data.Currency = data.Currency[:strings.Index(data.Currency, " ")]
			data.Currency = strings.Replace(data.Currency, "/", "", 1)
		}
		if debug {
			log.Println(print.PrettyPrint(data))
		}
		return data, nil
	}
	if strings.Index(msg, "TP1") == -1 || strings.Index(msg, "TP2") == -1 {
		err := errors.New("tp not found")
		return data, err
	}
	tab := strings.Split(msg, "\n")
	data.Trade = true
	for i := range tab {
		if strings.Index(tab[i], "/") > 0 {
			data.Currency = strings.Replace(tab[i], "/", "", 1)
			data.Currency = data.Currency[:strings.Index(data.Currency, "USDT")+len("USDT")]
			data.Currency = strings.Replace(data.Currency, " ", "", 1)
		}
		if strings.Index(tab[i], "nter ") > 0 {
			data.Entry = tab[i][strings.Index(tab[i], ": ")+2:]
		}
		if strings.Index(tab[i], "BUY") > 0 {
			data.Type = "Buy"
		} else if strings.Index(tab[i], "SELL") > 0 {
			data.Type = "Sell"
		}
		if strings.Index(tab[i], "TP1") > 0 {
			data.Tp1 = tab[i][strings.Index(tab[i], "- ")+2:]
		} else if strings.Index(tab[i], "TP2") > 0 {
			data.Tp2 = tab[i][strings.Index(tab[i], "- ")+2:]
		} else if strings.Index(tab[i], "TP3") > 0 {
			data.Tp3 = tab[i][strings.Index(tab[i], "- ")+2:]
		} else if strings.Index(tab[i], "SL") > 0 {
			data.Sl = tab[i][strings.Index(tab[i], "SL")+3:]
		} else if strings.Index(tab[i], "Leverage") > 0 {
			data.Level = tab[i][len("leverage ")+4 : strings.Index(tab[i], "x")]
		} else if strings.Index(tab[i], "set up as a") > 0 {
			data.SetUp = tab[i][strings.Index(tab[i], "set up as a")+len("set up as a ") : strings.Index(tab[i], " order")]
		} else if strings.Index(tab[i], "order or a") > 0 {
			data.Order = tab[i][strings.Index(tab[i], "order or a ")+len("order or a ") : strings.Index(tab[i], " order")]
		}
	}
	if debug {
		log.Println(print.PrettyPrint(data))
	}
	return data, nil
}
