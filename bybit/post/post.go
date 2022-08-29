package post

import (
	"bybit/bybit/bybit"
	"bybit/bybit/print"
	"bybit/bybit/sign"
	"bybit/env"
	"bytes"
	"encoding/json"
	"net/http"
)

func PostOrder(symbol string, api env.Env, trade *bybit.Trades) error {
	params := map[string]interface{}{
		"api_key":          api.Api,
		"side":             trade.GetType(symbol),
		"symbol":           symbol,
		"order_type":       "Limit",
		"price":            trade.GetSymbolPrice(symbol),
		"time_in_force":    "GoodTillCancel",
		"reduce_only":      false,
		"close_on_trigger": false,
		"stop_loss":        trade.GetSl(symbol),
	}
	// tp1
	_, err := sendPost(params, trade.GetTp1(symbol), api.Api_secret, trade, trade.GetTp1Order(symbol))
	if err != nil {
		return err
	}
	// tp2
	_, err = sendPost(params, trade.GetTp2(symbol), api.Api_secret, trade, trade.GetTp2Order(symbol))
	if err != nil {
		return err
	}
	// tp3
	_, err = sendPost(params, trade.GetTp3(symbol), api.Api_secret, trade, trade.GetTp3Order(symbol))
	if err != nil {
		return err
	}
	return nil
}

func sendPost(
	params map[string]interface{},
	tp string,
	api_secret string,
	trade *bybit.Trades,
	order string,
) (*http.Response, error) {
	var res Post

	params["take_profit"] = tp
	params["qty"] = order
	params["timestamp"] = print.GetTimestamp()
	params["sign"] = sign.GetSignedinter(params, api_secret)
	json_data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	println(print.PrettyPrint(params))
	req, err := http.Post(
		"https://api-testnet.bybit.com/private/linear/order/create",
		"application/json",
		bytes.NewBuffer(json_data),
	)
	if err != nil {
		return req, err
	}
	json.NewDecoder(req.Body).Decode(&res)
	trade.SetId(params["symbol"].(string), res.Result.OrderID)
	println(print.PrettyPrint(res))
	delete(params, "sign")
	delete(params, "take_profit")
	delete(params, "qty")
	delete(params, "timestamp")
	return req, nil
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
	req, err := http.Post(
		"https://api-testnet.bybit.com/private/linear/order/cancel-all",
		"application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		return err
	}
	json.NewDecoder(req.Body).Decode(&cancel)
	return nil
}

// func ChangeLs(api env.Env, params map[string]interface{}) error {
// 	params["sign"] = sign.GetSignedinter(params, api.Api_secret)
// 	json_data, err := json.Marshal(params)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = http.Post(
// 		"https://api-testnet.bybit.com/private/linear/position/trading-stop",
// 		"application/json",
// 		bytes.NewBuffer(json_data),
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	// json.NewDecoder(req.Body).Decode(&res)
// 	return nil
// }
