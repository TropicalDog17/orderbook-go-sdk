package types

import (
	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	decimal "github.com/shopspring/decimal"
)

type MarketSummary struct {
	MarketId string  `json:"marketId"`
	Open     float64 `json:"open"`
	High     float64 `json:"high"`
	Low      float64 `json:"low"`
	Volume   float64 `json:"volume"`
	Price    float64 `json:"price"`
	Change   float64 `json:"change"`
}

type SpotOrder struct {
	OrderType exchangetypes.OrderType `json:"orderType"`
	MarketId  string                  `json:"marketId"`
	Price     decimal.Decimal         `json:"price"`
	Quantity  decimal.Decimal         `json:"quantity"`
}

func NewSpotOrder(orderType exchangetypes.OrderType, marketId string, price decimal.Decimal, quantity decimal.Decimal) SpotOrder {
	return SpotOrder{
		OrderType: orderType,
		MarketId:  marketId,
		Price:     price,
		Quantity:  quantity,
	}
}
