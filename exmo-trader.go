package main

import (
	"exmo-trading/app/dataserver"
	"exmo-trading/app/trader"
	"exmo-trading/app/trader/strategies"
	config "exmo-trading/configs"
	"fmt"
	"path/filepath"
	"time"
)

type App struct {
	Traders      []trader.Trader
	DataHandlers []dataserver.Handler
}

func PrepareTrader(path, configName string, strategy strategies.Strategy) trader.Trader {
	trader := trader.Trader{}
	err := config.Load(&trader.Context, configName)
	if err != nil {
		fmt.Println(err)
	}
	strategy.Set(path + trader.Context.CandlesFile, trader.Context.CandlesFileVolume)
	trader.Context.TradesFile = path + trader.Context.TradesFile
	trader.Context.TradesHistoryFile = path + trader.Context.TradesHistoryFile
	trader.Context.CandlesFile = path + trader.Context.CandlesFile
	trader.Strategy = strategy
	return trader
}

func PrepareDataHandler(path, configName string) dataserver.Handler {
	datahandler := dataserver.Handler{}
	err := config.Load(&datahandler, configName)
	if err != nil {
		fmt.Println(err)
	}
	datahandler.CandlesFile = path + datahandler.CandlesFile
	return datahandler
}

func PrepareApp() *App {
	app := App{}
	fmt.Println("set up app...")
	time.Sleep(1 * time.Second)

	path, _ := filepath.Abs("./")

	fmt.Println("prepared traders: ")
	time.Sleep(1 * time.Second)
	app.Traders = append(app.Traders, PrepareTrader(path, path+"/configs/trader-configs/5min-btc-usdt-trader.yaml", &strategies.BollingerBandsTrader{}))
	app.Traders = append(app.Traders, PrepareTrader(path, path+"/configs/trader-configs/15min-btc-usdt-trader.yaml", &strategies.RSItrader{}))

	fmt.Println("prepared datahandlers: ")
	time.Sleep(1 * time.Second)
	app.DataHandlers = append(app.DataHandlers, PrepareDataHandler(path, path+"/configs/dataserver-configs/5min-btc-usdt-datahandler.yaml"))
	app.DataHandlers = append(app.DataHandlers, PrepareDataHandler(path, path+"/configs/dataserver-configs/15min-btc-usdt-datahandler.yaml"))

	return &app
}

func main() {
	App := PrepareApp()
	for i := range App.DataHandlers {
		fmt.Println("№" + fmt.Sprintf("%d", i+1) + " datahandler run ...")
		go App.DataHandlers[i].Run()
		time.Sleep(1 * time.Second)
	}
	time.Sleep(5 * time.Second)
	for i := range App.Traders {
		fmt.Println("№" + fmt.Sprintf("%d", i+1) + " trader run ...")
		go App.Traders[i].Run()
		time.Sleep(1 * time.Second)
	}
	for {
		time.Sleep(1000 * time.Second)
	}
}
