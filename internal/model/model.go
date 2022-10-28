package model

type SocketInstrument struct {
	Type string `json:"type"`
	Msg  FinancialInstrument
}

type SecurityIDs struct {
	Response []string
}

type NewInstrument struct {
	Response FinancialInstrument
}

type FinancialInstrument struct {
	SecurityID     string      `json:"securityID"`
	MDReqID        string      `json:"mdReqID"`
	Currency       string      `json:"currency"`
	Symbol         string      `json:"symbol"`
	Bid            []float32   `json:"Bid"`
	BidSize        []float32   `json:"BidSize"`
	Offer          []float32   `json:"Offer"`
	OfferSize      []float32   `json:"OfferSize"`
	Last           Last        `json:"last"`
	Underlying     any         `json:"underlying"`
	TradeVolume    TradeVolume `json:"tradeVolume"`
	SettlementType string      `json:"settlementType"`
}

type Last struct {
	Price float64 `json:"price"`
	Size  int     `json:"size"`
	Time  string  `json:"time"`
}

type TradeVolume struct {
	Price float64 `json:"price"`
	Size  int     `json:"size"`
	Time  string  `json:"time"`
}
