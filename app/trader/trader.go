package main

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/strategies"
	"exmo-trading/app/trader/traderutils"
	"fmt"
	"time"
)

type Event struct {
	Type    string
	Context *EventContext
}

type EventContext struct {
	Trend      string
	LastCandle *data.Candle
	OpenOrders *data.Orders
}

func (e *Event) Init(eventType string) error {
	parent := traderutils.GetCandlesAddr()
	trend := &strategies.GetTrend{}
	trend.Set(parent + "/cache/5min-candles.csv")
	OpenOrdersFile := parent + "/cache/open-orders.csv"
	var err error
	e.Context.Trend, err = trend.Analyze()
	if err != nil {
		return err
	}

	candles := &data.Candles{}
	candles.Array = make([]data.Candle, 0, 250)
	err = candles.Read(trend.CandlesFile)
	if err != nil {
		return err
	}

	e.Context.LastCandle = &candles.Array[len(candles.Array)-1]

	orders := &data.Orders{}
	orders.Array = make([]data.Order, 0, 10)
	err = orders.Read(OpenOrdersFile)
	if err != nil {
		return err
	}
	e.Context.OpenOrders = orders
	e.Type = eventType
	return nil
}

func (e *Event) Handle() {
	if len(e.Context.OpenOrders.Array) != 0 {
		lastOrder := &e.Context.OpenOrders.Array[len(e.Context.OpenOrders.Array)-1]
		if lastOrder.Action != e.Type {
			CloseOrder(lastOrder)
			newOrder := &data.Order{}
			newOrder.Action = e.Type
			OpenOrder(newOrder)
		}
	} else {
		newOrder := &data.Order{}
		newOrder.Action = e.Type
		OpenOrder(newOrder)
	}
}

func OpenOrder(o *data.Order) {

}

func CloseOrder(o *data.Order) {

}

type Trader interface {
	Analyze() (string, error)
	TestAnalyze() error
}

func Trade(t Trader) {
	eventType, err := t.Analyze()
	if err != nil {
		fmt.Println(err)
	}
	if eventType != "empty" {
		e := &Event{}
		err := e.Init(eventType)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func TestTrade(t Trader) {
	eventType := "empty"
	err := t.TestAnalyze()
	if err != nil {
		fmt.Println(err)
	}
	if eventType != "empty" {
		e := &Event{}
		err := e.Init(eventType)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	parent := traderutils.GetCandlesAddr()
	rsi := &strategies.GetTrend{}
	rsi.Set(parent + "/cache/5min-candles.csv")
	for {
		TestTrade(rsi)
		time.Sleep(60 * time.Second)
	}
}
