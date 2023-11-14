package main

import (
	"fmt"
	"github.com/maxzhovtyj/binance-test/pkg/binance"
	"github.com/sirupsen/logrus"
)

func main() {
	manager := binance.NewManager("", "")
	symbols, err := manager.GetExchangeInfo(5)
	if err != nil {
		logrus.Fatal(err)
	}

	list, err := manager.GetSymbolsPriceList(symbols)
	if err != nil {
		logrus.Fatal(err)
	}

	for _, res := range list {
		fmt.Println(res[binance.Symbol], res[binance.Price])
	}
}
