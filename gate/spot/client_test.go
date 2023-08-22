package spot

import (
	"context"
	"github.com/shopspring/decimal"
	"testing"
)

func TestClient_PriceWithAmount(t *testing.T) {
	tests := []struct {
		CurrencyPair string
		Amount       decimal.Decimal
	}{
		{"BTC_USDT", decimal.NewFromFloat(0.01)},
		{"ETH_USDT", decimal.NewFromFloat(0.01)},
	}
	client := New("", "")
	ctx := context.Background()
	for _, tt := range tests {
		buy, sell, err := client.PriceWithAmount(ctx, tt.CurrencyPair, tt.Amount)
		if err != nil {
			t.Errorf("get price error: %s", err)
			continue
		}
		t.Logf("%s, %s, %s, %s", tt.CurrencyPair, tt.Amount.String(), buy, sell)
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		Order order
		Price decimal.Decimal
	}{
		{
			Order: order{
				{
					Amount: decimal.NewFromInt(1),
					Price:  decimal.NewFromInt(1),
				},
			},
			Price: decimal.NewFromInt(1),
		},
		{
			Order: order{
				{
					Amount: decimal.NewFromInt(1),
					Price:  decimal.NewFromInt(1),
				},
				{
					Amount: decimal.NewFromInt(1),
					Price:  decimal.NewFromInt(2),
				},
			},
			Price: decimal.NewFromFloat(1.5),
		},
		{
			Order: order{
				{
					Amount: decimal.NewFromInt(1),
					Price:  decimal.NewFromInt(2),
				},
				{
					Amount: decimal.NewFromInt(1),
					Price:  decimal.NewFromInt(1),
				},
			},
			Price: decimal.NewFromFloat(1.5),
		},
		{
			Order: order{
				{
					Amount: decimal.NewFromInt(2),
					Price:  decimal.NewFromInt(1),
				},
				{
					Amount: decimal.NewFromInt(2),
					Price:  decimal.NewFromInt(2),
				},
				{
					Amount: decimal.NewFromInt(1),
					Price:  decimal.NewFromInt(4),
				},
			},
			Price: decimal.NewFromFloat(2),
		},
	}
	for i, tt := range tests {
		price := average(tt.Order)
		if !price.Equal(tt.Price) {
			t.Logf("case %d, expect %s, got %s", i, tt.Price, price)
		}
	}
}

func TestFill(t *testing.T) {
	tests := []struct {
		Orderbook [][]string
		Amount    decimal.Decimal
		Price     decimal.Decimal
	}{
		{
			Orderbook: [][]string{
				{"1", "1"},
				{"2", "2"},
			},
			Amount: decimal.NewFromFloat(1),
			Price:  decimal.NewFromFloat(1),
		},
		{
			Orderbook: [][]string{
				{"1", "1"},
				{"2", "2"},
			},
			Amount: decimal.NewFromFloat(2),
			Price:  decimal.NewFromFloat(1.5),
		},
		{
			Orderbook: [][]string{
				{"1", "1"},
				{"2", "2"},
				{"3", "10"},
			},
			Amount: decimal.NewFromFloat(3),
			Price:  decimal.NewFromFloat(1.6666666666666667),
		},
	}
	for i, tt := range tests {
		price := fill(tt.Amount, tt.Orderbook)
		if !price.Equal(tt.Price) {
			t.Errorf("case %d expect %s got %s", i, tt.Price, price)
		}
	}
}
