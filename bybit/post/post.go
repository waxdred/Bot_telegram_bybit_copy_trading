package post

import (
	"bybit/bybit/bybit"
	"bybit/bybit/get"
	"bybit/bybit/print"
	"bybit/bybit/sign"
	"bybit/env"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func PostOrder(symbol string, api env.Env, trade *bybit.Trades, debug bool) error {
	params := map[string]interface{}{
		"api_key":          api.Api,
		"side":             trade.GetType(symbol),
		"symbol":           symbol,
		"order_type":       "Limit",
		"price":            trade.GetEntry(symbol),
		"time_in_force":    "GoodTillCancel",
		"reduce_only":      false,
		"close_on_trigger": false,
		"stop_loss":        trade.GetSl(symbol),
	}
	// tp1
	_, err := sendPost(params, trade.GetTp1(symbol), api, trade, trade.GetTp1Order(symbol), debug)
	if err != nil {
		return err
	}
	// tp2
	_, err = sendPost(params, trade.GetTp2(symbol), api, trade, trade.GetTp2Order(symbol), debug)
	if err != nil {
		return err
	}
	// tp3
	_, err = sendPost(params, trade.GetTp3(symbol), api, trade, trade.GetTp3Order(symbol), debug)
	if err != nil {
		return err
	}
	return nil
}

func sendPost(
	params map[string]interface{},
	tp string,
	api env.Env,
	trade *bybit.Trades,
	order string,
	debug bool,
) (*http.Response, error) {
	var res Post

	params["take_profit"] = tp
	params["qty"] = order
	params["timestamp"] = print.GetTimestamp()
	params["sign"] = sign.GetSignedinter(params, api.Api_secret)
	json_data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	if debug {
		println(print.PrettyPrint(params))
	}
	url := fmt.Sprint(api.Url, "/private/linear/order/create")
	req, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return req, err
	}
	json.NewDecoder(req.Body).Decode(&res)
	if res.RetCode != 0 {
		return nil, errors.New(res.RetMsg)
	}
	if debug {
		log.Println(print.PrettyPrint(res))
	}
	trade.SetId(params["symbol"].(string), res.Result.OrderID)
	delete(params, "sign")
	delete(params, "take_profit")
	delete(params, "qty")
	delete(params, "timestamp")
	return req, nil
}

func PostIsoled(api env.Env, symbol string, trade *bybit.Trades, debug bool) error {
	var isolated Isolated
	params := map[string]interface{}{
		"api_key":       api.Api,
		"symbol":        symbol,
		"is_isolated":   true,
		"buy_leverage":  10,
		"sell_leverage": 10,
		"timestamp":     print.GetTimestamp(),
	}
	params["sign"] = sign.GetSignedinter(params, api.Api_secret)
	json_data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	url := fmt.Sprint(api.Url, "/private/linear/position/switch-isolated")
	req, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	json.NewDecoder(req.Body).Decode(&isolated)
	if debug {
		log.Printf("post PostIsoled")
		log.Println(print.PrettyPrint(isolated))
	}
	log.Printf("Isolated active: %d", params["buy_leverage"])
	return nil
}

func CancelOrder(symbol string, api env.Env, trade *bybit.Trades) error {
	params := map[string]string{
		"api_key": api.Api,
		"symbol":  symbol,
	}
	err := PostCancelOrder(params, api)
	if err != nil {
		return err
	}
	trade.Delete(symbol)
	log.Printf("Cancel order success: %s", symbol)
	return nil
}

func PostCancelOrder(params map[string]string, api env.Env) error {
	var cancel PostCancel

	params["timestamp"] = print.GetTimestamp()
	params["sign"] = sign.GetSigned(params, api.Api_secret)
	json_data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	delete(params, "sign")
	url := fmt.Sprint(api.Url, "/private/linear/order/cancel-all")
	req, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	json.NewDecoder(req.Body).Decode(&cancel)
	if cancel.RetCode != 0 {
		return errors.New(cancel.RetMsg)
	}
	return nil
}

func CancelBySl(price get.Price, trade *bybit.Trade) string {
	if trade.Type == "Buy" {
		val, _ := strconv.ParseFloat(price.Result[0].BidPrice, 8)
		val = (val * 0.01 / 100) + val
		return fmt.Sprint("%4.v", val)
	} else if trade.Type == "Sell" {
		val, _ := strconv.ParseFloat(price.Result[0].BidPrice, 8)
		val = (val * 0.01 / 100) - val
		return fmt.Sprint("%4.v", val)
	}
	return ""
}

func ChangeLs(api env.Env, symbol string, sl string) error {
	params := map[string]string{
		"api_key":   api.Api,
		"symbol":    symbol,
		"stop_loss": sl,
		"timestamp": print.GetTimestamp(),
	}
	params["sign"] = sign.GetSigned(params, api.Api_secret)
	json_data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	url := fmt.Sprint(api.Url, "/private/linear/position/trading-stop")
	_, err = http.Post(
		url,
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return err
	}
	log.Printf("Change ls: %s", symbol)
	return nil
}
