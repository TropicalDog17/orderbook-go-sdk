package exchange

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	exchangeclient "github.com/InjectiveLabs/sdk-go/client/exchange"
	"github.com/TropicalDog17/orderbook-go-sdk/internal/chain"
	types "github.com/TropicalDog17/orderbook-go-sdk/internal/types"
)

var _ CronosClient = (*MbClient)(nil)
var _ ExchangeFetcher = (*MbClient)(nil)

type ExchangeFetcher interface {
	GetPrice(ticker string) (float64, error)
}

type WalletFetcher interface {
	GetBalance() (float64, error)
}

type CronosClient interface {
	GetMarketSummary(marketId string) (types.MarketSummary, error)
}
type MbClient struct {
	exchangeClient exchangeclient.ExchangeClient
	chainClient    *chain.ChainClient
	config         *types.Config
}

func NewMbClient(networkType string, config *types.Config) *MbClient {
	if networkType != "local" {
		panic("Only local network type is supported")
	}

	network := types.DefaultNetwork()
	exchangeClient, err := exchangeclient.NewExchangeClient(network)
	if err != nil {
		panic(err)
	}
	chainClient := chain.NewChainClient("genesis") // TODO: refactor hard code
	return &MbClient{
		exchangeClient: exchangeClient,
		chainClient:    &chainClient,
		config:         config,
	}
}

func (c *MbClient) GetPrice(ticker string) (float64, error) {
	ticker = strings.Replace(ticker, "-", "", -1)
	ticker = strings.Replace(ticker, "/", "", -1)
	ticker = strings.ToUpper(ticker)
	marketId := os.Getenv(ticker)
	if marketId == "" {
		return 0, fmt.Errorf("marketId not found for ticker %s", ticker)
	}
	marketSummary, err := c.GetMarketSummary(marketId)
	if err != nil {
		return 0, err
	}
	return marketSummary.Price, nil
}

func (c *MbClient) GetMarketSummary(marketId string) (types.MarketSummary, error) {
	// TODO: fix hard code

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	endpoint := fmt.Sprintf("%s/api/chronos/v1/spot/market_summary?marketId=%s&resolution=24h", c.config.ChronosEndpoint, marketId)
	var marketSummary types.MarketSummary
	resp, err := client.Get(endpoint)

	if err != nil {
		return marketSummary, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return marketSummary, err
	}
	if err := json.Unmarshal(bodyBytes, &marketSummary); err != nil {
		return marketSummary, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return marketSummary, nil
}

// func (c *MbClient) GetMarketsAssistant() chainclient.MarketsAssistant {
// 	ctx := context.Background()

// 	marketsAssistant, err := chainclient.NewMarketsAssistantInitializedFromChain(ctx, *c.exchangeClient)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return marketsAssistant
// }

func (c *MbClient) GetChainClient() *chain.ChainClient {
	return c.chainClient
}

func (c *MbClient) GetDecimals(ctx context.Context, marketId string) (baseDecimal, quoteDecimal int32) {
	market, err := c.exchangeClient.GetSpotMarket(ctx, marketId)
	if err != nil {
		panic(err)
	}
	baseDecimal = market.Market.BaseTokenMeta.Decimals
	quoteDecimal = market.Market.QuoteTokenMeta.Decimals
	return baseDecimal, quoteDecimal
}
