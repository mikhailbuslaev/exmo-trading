package main

import (
	"fmt"
	"time"
	"exmo-trading/app/trader"
	"exmo-trading/app/dataserver"
	"exmo-trading/app/trader/strategies"
	"exmo-trading/configs"
	"path/filepath"
)

type App struct {
	Traders []trader.Trader
	DataHandlers []dataserver.Handler
}

func PrepareTrader(path, configName string, analyzer trader.Analyzer) trader.Trader {
	trader := trader.Trader{}
	trader.Analyzer = analyzer
	err := config.Load(&trader.Context, configName)
	if err != nil {
		fmt.Println(err)
	}
	trader.Context.TradesFile = path+trader.Context.TradesFile
	trader.Context.TradesHistoryFile = path+trader.Context.TradesHistoryFile
	trader.Context.CandlesFile  = path+trader.Context.CandlesFile
	return trader
}

func PrepareDataHandler(path, configName string) dataserver.Handler {
	datahandler := dataserver.Handler{}
	err := config.Load(&datahandler, configName)
	if err != nil {
		fmt.Println(err)
	}
	datahandler.CandlesFile  = path+datahandler.CandlesFile
	return datahandler
}

func PrepareApp() *App{

	fmt.Println("set up app...")
	time.Sleep(1*time.Second)

	path, _ := filepath.Abs("./")

	rsi_5min := strategies.RSItrader{}
	rsi_5min.Set(path+"/cache/5min-btc-usdt-candles.csv")
	rsi_15min := strategies.RSItrader{}
	rsi_15min.Set(path+"/cache/15min-btc-usdt-candles.csv")
	traders := make([]trader.Trader, 0, 10)
	datahandlers := make([]dataserver.Handler, 0, 10)

	fmt.Println("prepared traders: ")
	time.Sleep(1*time.Second)
	traders = append(traders, PrepareTrader(path, path+"/configs/trader-configs/5min-btc-usdt-trader.yaml", &rsi_5min))
	traders = append(traders, PrepareTrader(path, path+"/configs/trader-configs/15min-btc-usdt-trader.yaml", &rsi_15min))

	fmt.Println("prepared datahandlers: ")
	time.Sleep(1*time.Second)
	datahandlers = append(datahandlers, PrepareDataHandler(path, path+"/configs/dataserver-configs/5min-btc-usdt-datahandler.yaml"))
	datahandlers = append(datahandlers, PrepareDataHandler(path, path+"/configs/dataserver-configs/15min-btc-usdt-datahandler.yaml"))
	
	return &App{Traders: traders, DataHandlers: datahandlers}
}

func main() {
	App := PrepareApp()
	for i := range App.DataHandlers {
		fmt.Println("№"+fmt.Sprintf("%d", i+1)+" datahandler run ...")
		go App.DataHandlers[i].Run()
		time.Sleep(1*time.Second)
	}
	time.Sleep(5*time.Second)
	for i := range App.Traders {
		fmt.Println("№"+fmt.Sprintf("%d", i+1)+" trader run ...")
		go App.Traders[i].Run()
		time.Sleep(1*time.Second)
	}
	for{time.Sleep(5*time.Second)}
}
