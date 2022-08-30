package bybit

import (
	"bybit/bybit/get"
	"bybit/bybit/print"
	"bybit/bybit/telegram"
	"bybit/env"
	"fmt"
	"log"
	"math"
	"strconv"
)

type Trades []Trade

func RoundFloat(val float64, precision uint) string {
	ratio := math.Pow(10, float64(precision))
	ret := math.Round(val*ratio) / ratio
	rets := fmt.Sprint(ret)
	return rets
}

func (t *Trades) SetId(symbol string, id string) {
	ls := *t

	for i := 0; i < len(ls); i++ {
		if ls[i].Symbol == symbol {
			ls[i].Id = append(ls[i].Id, id)
		}
	}
}

func (t *Trades) GetSymbolOrder() []string {
	ls := *t
	var ret []string
	var check bool

	for i := 0; i < len(ls); i++ {
		check = true
		if ret == nil {
			ret = append(ret, ls[i].Symbol)
		} else {
			for j := 0; j < len(ret); j++ {
				if ret[j] == ls[i].Symbol {
					check = false
				}
			}
			if check {
				ret = append(ret, ls[i].Symbol)
			}
		}
	}
	return ret
}

func (t *Trades) CheckSymbol(symbol string) bool {
	ls := *t

	for i := 0; i < len(ls); i++ {
		if ls[i].Symbol == symbol {
			return false
		}
	}
	return true
}

func GetTrade(symbol string, t *Trades) *Trade {
	ls := *t

	for i := 0; i < len(ls); i++ {
		if ls[i].Symbol == symbol {
			return &ls[i]
		}
	}
	return nil
}

func (t *Trades) Add(api env.Env, data telegram.Data, price get.Price) bool {
	wallet := get.GetWallet(api)
	available := wallet.Result.Usdt.AvailableBalance / 3
	prices, _ := strconv.ParseFloat(price.Result[0].BidPrice, 8)
	log.Println(print.PrettyPrint(available))
	elem := Trade{
		Symbol:      data.Currency,
		Type:        data.Type,
		Order:       data.Order,
		SymbolPrice: price.Result[0].BidPrice,
		Wallet:      fmt.Sprint(RoundFloat(available, 4)),
		Entry:       data.Entry,
		Tp1Order:    RoundFloat((available*50/100)/prices, 4),
		Tp2Order:    RoundFloat((available*25/100)/prices, 4),
		Tp3Order:    RoundFloat((available*15/100)/prices, 4),
		Tp1:         data.Tp1,
		Tp2:         data.Tp2,
		Tp3:         data.Tp3,
		Sl:          data.Sl,
	}
	log.Printf("available")
	if t.CheckSymbol(data.Currency) == false {
		log.Printf("Trade actif Symbol: %s", data.Currency)
		return false
	}
	*t = append(*t, elem)
	return true
}

func (t *Trades) Delete(symbol string) bool {
	tmp := &Trades{}
	check := 0
	for i := 0; i < len(*t); i++ {
		if (*t)[i].Symbol != symbol {
			*tmp = append(*tmp, (*t)[i])
		} else {
			check = 1
		}
	}
	*t = *tmp
	if check == 1 {
		return false
	}
	return true
}

func (t *Trades) Print() {
	ls := *t

	for i := 0; i < len(ls); i++ {
		log.Println(print.PrettyPrint(ls[i]))
	}
}

// get Trade
func (t *Trades) GetSymbol(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Symbol
	}
	return ""
}

func (t *Trades) GetType(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Type
	}
	return ""
}

func (t *Trades) GetOrder(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Order
	}
	return ""
}

func (t *Trades) GetSymbolPrice(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.SymbolPrice
	}
	return ""
}

func (t *Trades) GetWallet(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Wallet
	}
	return ""
}

func (t *Trades) GetPrice(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Price
	}
	return ""
}

func (t *Trades) GetEntry(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Entry
	}
	return ""
}

func (t *Trades) GetTp1Order(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Tp1Order
	}
	return ""
}

func (t *Trades) GetTp2Order(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Tp2Order
	}
	return ""
}

func (t *Trades) GetTp3Order(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Tp3Order
	}
	return ""
}

func (t *Trades) GetTp1(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Tp1
	}
	return ""
}

func (t *Trades) GetTp2(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Tp2
	}
	return ""
}

func (t *Trades) GetTp3(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Tp2
	}
	return ""
}

func (t *Trades) GetSl(symbol string) string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Sl
	}
	return ""
}

func (t *Trades) GetId(symbol string) []string {
	ret := GetTrade(symbol, t)
	if ret != nil {
		return ret.Id
	}
	return nil
}
