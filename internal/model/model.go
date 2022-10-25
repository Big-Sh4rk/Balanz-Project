package model

type Response struct {
	AL300002CCTARS  FinancialInstrument `json:"AL30-0002-C-CT-ARS"`
	GD30D0001CCTUSD FinancialInstrument `json:"GD30D-0001-C-CT-USD"`
	AL300001CCTARS  FinancialInstrument `json:"AL30-0001-C-CT-ARS"`
	AL300003CCTARS  FinancialInstrument `json:"AL30-0003-C-CT-ARS"`
	GD300001CCTARS  FinancialInstrument `json:"GD30-0001-C-CT-ARS"`
	AL30D0003CCTUSD FinancialInstrument `json:"AL30D-0003-C-CT-USD"`
	GD300003CCTARS  FinancialInstrument `json:"GD30-0003-C-CT-ARS"`
	GD30D0003CCTUSD FinancialInstrument `json:"GD30D-0003-C-CT-USD"`
	AL30D0001CCTUSD FinancialInstrument `json:"AL30D-0001-C-CT-USD"`
	AL30C0003CCTEXT FinancialInstrument `json:"AL30C-0003-C-CT-EXT"`
	GD30C0003CCTEXT FinancialInstrument `json:"GD30C-0003-C-CT-EXT"`
	GD30C0001CCTEXT FinancialInstrument `json:"GD30C-0001-C-CT-EXT"`
	AL30C0001CCTEXT FinancialInstrument `json:"AL30C-0001-C-CT-EXT"`
	AL30D0002CCTUSD FinancialInstrument `json:"AL30D-0002-C-CT-USD"`
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
	Price float32 `json:"price"`
	Size  int     `json:"size"`
	Time  string  `json:"time"`
}

type TradeVolume struct {
	Price float32 `json:"price"`
	Size  int     `json:"size"`
	Time  string  `json:"time"`
}
