package main

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/strategies"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Trader interface {
	Analyze() (*data.Order, error)
}

func Trade(t Trader, ordersFile string) {
	order, err := t.Analyze()
	if err != nil {
		fmt.Println(err)
	}
	if order != nil {
		ExecOrder(order, ordersFile)
	} else {
		fmt.Println("no new orders")
	}
}

func ExecOrder(o *data.Order, ordersFile string) {
	o.Print()
	o.Write(ordersFile)
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	parent := filepath.Dir(filepath.Dir(wd))
	rsi := &strategies.RSItrader{}
	rsi.Set(parent + "/cache/5min-candles.csv")
	for {
		Trade(rsi, parent+"/cache/5min-orders.csv")
		time.Sleep(30 * time.Second)
	}
}
