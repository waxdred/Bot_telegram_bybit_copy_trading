package bybit

type Trade struct {
	Symbol      string `json:"symbol"`
	Type        string `json:"type"`
	Order       string `json:"order"`
	SymbolPrice string `json:"symbolPrice"`
	Wallet      string `json:"wallet"`
	Price       string `json:"price"`
	Entry       string `json:"entry"`
	Tp1Order    string `json:"tp_1Order"`
	Tp2Order    string `json:"tp_2Order"`
	Tp3Order    string `json:"tp_3Order"`
	Tp1         string `json:"tp1"`
	Tp2         string `json:"tp2"`
	Tp3         string `json:"tp3"`
	Sl          string `json:"Sl"`
	Id          []string
}
