package main

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/strategies"
	"fmt"
	"time"
)

type Trader interface {
	Analyze() (*data.Order, error)
}

func Trade(t Trader) {
	order, err := t.Analyze()
	if err != nil {
		fmt.Println(err)
	}
	if order != nil {
		ExecOrder(order)
	} else {
		fmt.Println("no new orders")
	}
}

func ExecOrder(o *data.Order) {
	o.Print()
	o.Write("cache/5min-orders.csv")
}

func main() {
	ma := &strategies.MAintersectionTrader{}
	ma.Set()
	rsi := &strategies.RSItrader{}
	rsi.Set()
	for {
		Trade(ma)
		Trade(rsi)
		time.Sleep(30 * time.Second)
	}
}
