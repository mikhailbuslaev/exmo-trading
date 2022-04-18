package main

import (
	"fmt"
	"log"
	"os/exec"
	"exmo-trading/app/trader"
	"exmo-trading/app/dataserver"
	"exmo-trading/app/strategies"
	"exmo-trading/configs"
)

type App struct {
	Traders []trader.Trader
	DataHandlers []dataserver.Handler
}

func PrepareTrader(configName string, analyzer Analyzer) trader.Trader {
	trader := trader.Trader{}
	trader.Analyzer = analyzer
	err := config.Load(configName, &trader)
	if err != nil {
		fmt.Println(err)
	}
	return trader
}

func PrepareDataHandler(configName string) dataserver.Handler {
	datahandler := dataserver.Handler{}
	err := config.Load(configName, &datahandler)
	if err != nil {
		fmt.Println(err)
	}
	return datahandler
}

func PrepareApp() *App{
	rsi := strategies.RSItrader{}
	rsi.Set()
	traders := make([]trader.Trader, 0, 10)
	datahandlers := make([]dataserver.Handler, 0, 10)
	traders = append(traders, PrepareTrader("/configs/trader-configs/5min-btc-usdt-trader.yaml", rsi))
	traders := append(traders, PrepareTrader("/configs/trader-configs/15min-btc-usdt-trader.yaml", rsi))
	datahandlers := append(datahandlers, PrepareDataHandler("/configs/dataserver-configs/5min-btc-usdt-datahandler.yaml", rsi))
	datahandlers := append(datahandlers, PrepareDataHandler("/configs/dataserver-configs/15min-btc-usdt-datahandler.yaml", rsi))
	return &App{Traders: traders, DataHandlers: datahandlers}
}

func main() {
	App := PrepareApp()

	for i := range App.Traders {
		go App.Traders[i].Run
	}

	for i := range App.DataHandlers {
		go App.DataHandlers[i].Run
	}
}
