package binance

import (
	"context"
	"fmt"
	"github.com/aiviaio/go-binance/v2"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	Symbol = "symbol"
	Price  = "price"
)

type ListPrice map[string]string

type Manager struct {
	client *binance.Client
}

func NewManager(apiKey, secretKey string) *Manager {
	return &Manager{client: binance.NewClient(apiKey, secretKey)}
}

func (m *Manager) GetExchangeInfo(limit int) ([]string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	service, err := m.client.NewExchangeInfoService().Do(ctx, []binance.RequestOption{}...)
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange info, %v", err)
	}

	var symbols []string

	if limit > 0 && limit < len(service.Symbols) {
		service.Symbols = service.Symbols[:limit]
	}

	for _, s := range service.Symbols {
		symbols = append(symbols, s.Symbol)
	}

	return symbols, nil
}

func (m *Manager) GetSymbolsPriceList(symbols []string) ([]ListPrice, error) {
	ch := make(chan ListPrice, len(symbols))
	done := make(chan struct{})

	for _, symbol := range symbols {
		go func(res chan<- ListPrice, s string, done chan<- struct{}) {
			symbolPrice, err := m.GetSymbolPrice(s)
			if err != nil {
				logrus.Errorf("failed to load '%s', %v", s, err)
				res <- ListPrice{}
				return
			}

			res <- symbolPrice
			done <- struct{}{}
		}(ch, symbol, done)
	}

	go func() {
		for i := 0; i < len(symbols); i++ {
			<-done
		}
		close(ch)
	}()

	var res []ListPrice

	for price := range ch {
		res = append(res, price)
	}

	return res, nil
}

func (m *Manager) GetSymbolPrice(symbol string) (ListPrice, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	res, err := m.client.NewListPricesService().Symbol(symbol).Do(ctx, []binance.RequestOption{}...)
	if err != nil {
		return nil, fmt.Errorf("failed to get list prices, symbol '%s', %v", symbol, err)
	}

	if len(res) < 1 {
		return nil, err
	}

	return ListPrice{
		Symbol: res[0].Symbol,
		Price:  res[0].Price,
	}, nil
}
