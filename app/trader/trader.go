package main

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/strategies"
	"exmo-trading/app/trader/traderutils"
	"fmt"
	"time"
)

type Event struct {
	Action    string
	Context *EventContext
}

type EventContext struct {
	Trend      string
	LastCandle *data.Candle
	OrdersFile string
	OrdersHistoryFile string
	Orders *data.Orders
}

func (e *Event) Init(eventType string) error {
	parent := traderutils.GetCandlesAddr()
	trend := &strategies.GetTrend{}
	trend.Set(parent + "/cache/5min-candles.csv")
	e.OrdersFile = parent + "/cache/orders.csv"
	e.OrdersHistoryFile = parent + "/cache/orders-history.csv"
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
	err = orders.Read(e.OrdersFile)
	if err != nil {
		return err
	}
	e.Context.Orders = orders
	e.Action = eventType
	return nil
}

func (e *Event) Handle() error{
	if len(e.Context.Orders.Array) != 0 && e.Context.Orders.Array[len(e.Context.Orders.Array)-1].Action != e.Action{
		lastOrder := &e.Context.Orders.Array[len(e.Context.Orders.Array)-1]
		if lastOrder.Action != e.Action {
			err := e.CloseOrder(lastOrder)
			if err != nil {
				return err
			}
			err = e.OpenOrder()
			if err != nil {
				return err
			}
		}
	} else if len(e.Context.Orders.Array) == 0{	
		err != e.OpenOrder(newOrder)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Event) OpenOrder() error {
	newOrder := &data.Order{}
	newOrder.Action = e.Type
	now := time.Now()
	newOrder.Id = now.Unix()
	if e.Action == "long" {
		newOrder.StopLimit = e.Context.LastCandle.Close - e.Context.LastCandle.Close*0.01
	} else if e.Action == "short"{
		newOrder.StopLimit = e.Context.LastCandle.Close + e.Context.LastCandle.Close*0.01
	}
	newOrder.Price = e.Context.LastCandle.Close
	newOrder.Status = "open"
	e.Context.Orders.Array[len(e.Context.Orders.Array)] = newOrder
	err := Rewrite(e.Context.Orders, e.OrdersFile)
	if err != nil {
		return err
	}
}

func (e *Event) CloseOrder(o *data.Order) error{
	var orderIndex int
	for i := range e.Context.Orders {
		if o.Id == e.Context.Orders.Array[i].Id {
			orderIndex = i
		}
	}
	e.Context.Orders.Array[orderIndex].Status = "close"
	err := &e.Context.Orders.Array[orderIndex].Write(e.Context.OrdersHistoryFile)
	if err != nil {
		return err
	}
	e.Context.Orders.Array = append(e.Context.Orders.Array[:orderIndex], a[orderIndex+1:]...)

	err = Rewrite(e.Context.Orders, e.OrdersFile)
	if err != nil {
		return err
	}
	return nil
}

type Trader interface {
	Analyze() (string, error)
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
		err = e.Handle()
		if err != nil {
			fmt.Println(err)
		}
	}
}


func main() {
	parent := traderutils.GetParentDir()
	rsi := &strategies.RSItrader{}
	rsi.Set(parent + "/cache/5min-candles.csv")
	for {
		Trade(rsi)
		time.Sleep(60 * time.Second)
	}
}
