package exchange

import (
	"context"
	"fmt"

	exchangetypes "github.com/InjectiveLabs/sdk-go/chain/exchange/types"
	"github.com/TropicalDog17/orderbook-go-sdk/internal/types"
	utils "github.com/TropicalDog17/orderbook-go-sdk/pkg/utils"
	"github.com/google/uuid"
)

type OrderMaker interface {
	PlaceSpotOrder(order types.SpotOrder) error
	PlaceMarketOrder() error
}

func (c *MbClient) PlaceSpotOrder(order types.SpotOrder) (string, error) {
	chainClient := c.chainClient.GetInjectiveChainClient()
	senderAddress := c.chainClient.SenderAddress
	ctx := context.Background()

	defaultSubaccountID := chainClient.DefaultSubaccount(senderAddress)
	baseDecimal, quoteDecimal := c.GetDecimals(ctx, order.MarketId)
	spotOrder := exchangetypes.SpotOrder{
		OrderType: exchangetypes.OrderType_BUY,
		MarketId:  order.MarketId,
		OrderInfo: exchangetypes.OrderInfo{
			SubaccountId: defaultSubaccountID.String(),
			Price:        utils.PriceToChainFormat(order.Price, baseDecimal, quoteDecimal),
			Quantity:     utils.QuantityToChainFormat(order.Quantity, baseDecimal),
			Cid:          uuid.NewString(),
		},
	}
	fmt.Println("spot order: ", spotOrder)
	msg := new(exchangetypes.MsgCreateSpotLimitOrder)
	msg.Sender = senderAddress.String()
	msg.Order = spotOrder
	simRes, err := chainClient.SimulateMsg(chainClient.ClientContext(), msg)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	msgCreateSpotLimitOrderResponse := exchangetypes.MsgCreateSpotLimitOrderResponse{}
	err = msgCreateSpotLimitOrderResponse.Unmarshal(simRes.Result.MsgResponses[0].Value)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//AsyncBroadcastMsg, SyncBroadcastMsg, QueueBroadcastMsg

	tx, err := chainClient.SyncBroadcastMsg(msg)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	txHash := tx.TxResponse.TxHash

	return txHash, nil
}

// func (c *MbClient) PlaceLimitOrder(order types.LimitOrder) error {
// 	chainClient := c.chainClient.GetInjectiveChainClient()
// 	senderAddress := c.chainClient.SenderAddress
// 	ctx := context.Background()

// 	defaultSubaccountID := chainClient.DefaultSubaccount(senderAddress)
// 	baseDecimal, quoteDecimal := c.GetDecimals(ctx, order.MarketId)
// 	spotOrder := exchangetypes.SpotOrder{
// 		OrderType: exchangetypes.OrderType_BUY,
// 		MarketId:  "0xfbd55f13641acbb6e69d7b59eb335dabe2ecbfea136082ce2eedaba8a0c917a3",
// 		OrderInfo: exchangetypes.OrderInfo{
// 			SubaccountId: defaultSubaccountID.String(),
// 			Price:        utils.PriceToChainFormat(order.Price, baseDecimal, quoteDecimal),
// 			Quantity:     utils.QuantityToChainFormat(order.Quantity, baseDecimal),
// 			Cid:          uuid.NewString(),
// 		},
// 	}
// 	fmt.Println("spot order: ", spotOrder)
// 	msg := new(exchangetypes.MsgCreateSpotLimitOrder)
// 	msg.Sender = senderAddress.String()
// 	msg.Order = spotOrder
// 	simRes, err := chainClient.SimulateMsg(chainClient.ClientContext(), msg)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	msgCreateSpotLimitOrderResponse := exchangetypes.MsgCreateSpotLimitOrderResponse{}
// 	err = msgCreateSpotLimitOrderResponse.Unmarshal(simRes.Result.MsgResponses[0].Value)

// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	//AsyncBroadcastMsg, SyncBroadcastMsg, QueueBroadcastMsg

// 	tx, err := chainClient.SyncBroadcastMsg(msg)

// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	fmt.Println("tx: ", tx)

// 	return nil
// }
