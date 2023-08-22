package spot

import (
	"context"
	"github.com/gateio/gateapi-go/v6"
	"github.com/shopspring/decimal"
)

type Client struct {
	key    string
	secret string
}

func New(key, secret string) *Client {
	return &Client{
		key:    key,
		secret: secret,
	}
}

func (c *Client) PriceWithAmount(ctx context.Context, currencyPair string, amount decimal.Decimal) (buy, sell decimal.Decimal, err error) {
	client := gateapi.NewAPIClient(gateapi.NewConfiguration())
	//auth := context.WithValue(context.Background(), gateapi.ContextGateAPIV4, gateapi.GateAPIV4{
	//	Key:    "YOUR_API_KEY",
	//	Secret: "YOUR_API_SECRET",
	//})
	//currencyPair := "BTC_USDT" // string - Currency pair

	result, _, err := client.SpotApi.ListOrderBook(ctx, currencyPair, nil)
	if err != nil {
		return
	}
	// ask 是卖方深度列表，价格从低到高
	// bid 是买方深度列表，价格从高到低
	buy = fill(amount, result.Asks)
	sell = fill(amount, result.Bids)

	return
}

type trade struct {
	Amount decimal.Decimal
	Price  decimal.Decimal
}

type order []trade

func fill(amount decimal.Decimal, ask [][]string) decimal.Decimal {
	var o order
	left := amount
	for _, a := range ask {
		// 0: 价格,  1: 数量
		canTake := decimal.RequireFromString(a[1])
		price := decimal.RequireFromString(a[0])
		t := trade{
			Price: price,
		}
		if canTake.GreaterThanOrEqual(left) {
			t.Amount = left
		} else {
			t.Amount = canTake
		}
		o = append(o, t)
		left = left.Sub(t.Amount)
		if left.IsZero() {
			break
		}
	}
	return average(o)
}

func average(o order) decimal.Decimal {
	amount := decimal.NewFromInt(0)
	total := decimal.NewFromInt(0)
	for _, t := range o {
		amount = amount.Add(t.Amount)
		total = total.Add(t.Amount.Mul(t.Price))
	}
	if amount.IsZero() {
		return amount
	}
	return total.Div(amount)
}
