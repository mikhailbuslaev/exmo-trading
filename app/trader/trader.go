package trader

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

func Launch() {
	m := &strategies.MAintersectionTrader{}
	m.Set()
	for {
		Trade(m)
		time.Sleep(30 * time.Second)
	}
}
