package main // trader is autonomic microservice that launch strategies and handling signals by making events

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/signals"
	"exmo-trading/app/trader/strategies"
	"exmo-trading/app/trader/traderutils"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	TradesFile        string  = "/cache/trades.csv"
	TradesHistoryFile string  = "/cache/trades-history.csv"
	CandlesFile       string  = "/cache/5min-candles.csv"
	CandlesFileVolume int     = 250  // volume of candles in candles file
	TradesFileVolume  int     = 10   // max amount of open trades
	StopLimitPercent  float64 = 0.01 // 0.01 means 1%, when price change goes higher or lower this percent, we close unprofitable trade
	TradeAmount float64 = 100 // amount of usdt in trade
	TraderTimeout time.Duration = 30 // time of trader reload
)

type Event struct { // event is need to making trades from strategy signals, we making trade relying on context
	Action  string
	Context EventContext
}

type EventContext struct {
	LastPrice         float64
	TradesFile        string
	TradesHistoryFile string
	Trades            data.Trades
}

type Trader interface {
	Analyze() (string, error)
}

func GetLastPrice() (float64, error) {
	url := "https://api.exmo.me/v1.1/ticker"
	method := "POST"
	payload := strings.NewReader("")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	ticker := &data.Ticker{}

	err = ticker.ParseJsonTickers([]byte(body), "BTC_USDT")
	if err != nil {
		return 0, err
	}
	price, err := strconv.ParseFloat(ticker.Avg, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (e *Event) Init(eventType string) error { // loading context for our event
	parent := traderutils.GetParentDir()
	e.Context.TradesFile = parent + TradesFile
	e.Context.TradesHistoryFile = parent + TradesHistoryFile

	lastPrice, err := GetLastPrice()
	if err != nil {
		return err
	}
	e.Context.LastPrice = lastPrice

	trades := &data.Trades{}
	trades.Array = make([]data.Trade, 0, TradesFileVolume)
	err = trades.Read(e.Context.TradesFile)
	if err != nil {
		return err
	}
	e.Context.Trades = *trades
	e.Action = eventType
	return nil
}

func (e *Event) HandleOpenedTrades() error {
	for i := range e.Context.Trades.Array {
		if e.Context.Trades.Array[i].Action == signals.Long && e.Context.Trades.Array[i].StopLimit > e.Context.LastPrice {
			err := e.CloseTrade(&e.Context.Trades.Array[i])
			if err != nil {
				return err
			}
		}
		if e.Context.Trades.Array[i].Action == signals.Short && e.Context.Trades.Array[i].StopLimit < e.Context.LastPrice {
			err := e.CloseTrade(&e.Context.Trades.Array[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Event) HandleSignal() error {
	if e.Action == signals.Long || e.Action == signals.Short {
		if len(e.Context.Trades.Array) != 0 && e.Context.Trades.Array[len(e.Context.Trades.Array)-1].Action != e.Action {
			lastTrade := &e.Context.Trades.Array[len(e.Context.Trades.Array)-1]
			if lastTrade.Action != e.Action {
				err := e.CloseTrade(lastTrade)
				if err != nil {
					return err
				}
				err = e.OpenTrade()
				if err != nil {
					return err
				}
			}
		}
		if len(e.Context.Trades.Array) == 0 {
			err := e.OpenTrade()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Event) Handle() error { // checking event context and making trades
	err := e.HandleOpenedTrades()
	if err != nil {
		return err
	}
	err = e.HandleSignal()
	if err != nil {
		return err
	}
	return nil
}

func (e *Event) OpenTrade() error {
	newTrade := &data.Trade{}
	newTrade.Action = e.Action
	now := time.Now()
	newTrade.Id = now.Unix()
	if e.Action == signals.Long {
		newTrade.StopLimit = e.Context.LastPrice - e.Context.LastPrice*StopLimitPercent
	} else if e.Action == signals.Short {
		newTrade.StopLimit = e.Context.LastPrice + e.Context.LastPrice*StopLimitPercent
	}
	newTrade.OpenPrice = e.Context.LastPrice
	newTrade.Status = signals.TradeOpened
	newTrade.Quantity = TradeAmount
	e.Context.Trades.Array = append(e.Context.Trades.Array, *newTrade)
	err := data.Rewrite(&e.Context.Trades, e.Context.TradesFile)
	if err != nil {
		return err
	}
	fmt.Println("trade opened")
	return nil
}

func (e *Event) CloseTrade(t *data.Trade) error {
	var tradeIndex int
	for i := range e.Context.Trades.Array {
		if t.Id == e.Context.Trades.Array[i].Id {
			tradeIndex = i
		}
	}
	e.Context.Trades.Array[tradeIndex].Status = signals.TradeClosed
	lastPrice, err := GetLastPrice()
	if err != nil {
		return err
	}
	e.Context.Trades.Array[tradeIndex].ClosePrice = lastPrice
	err = e.Context.Trades.Array[tradeIndex].Write(e.Context.TradesHistoryFile) // write trade to archive
	if err != nil {
		return err
	}
	e.Context.Trades.Array = append(e.Context.Trades.Array[:tradeIndex], e.Context.Trades.Array[tradeIndex+1:]...) // deleting trade

	err = data.Rewrite(&e.Context.Trades, e.Context.TradesFile) // rewrite open trades file
	if err != nil {
		return err
	}
	fmt.Println("trade closed")
	return nil
}

func Trade(t Trader) {
	eventType, err := t.Analyze()
	if err != nil {
		fmt.Println(err)
	}

	e := &Event{}
	err = e.Init(eventType)
	if err != nil {
		fmt.Println(err)
	}

	err = e.Handle()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("launching trader...")
	parent := traderutils.GetParentDir()
	rsi := &strategies.RSItrader{}
	rsi.Set(parent + CandlesFile)
	for {
		Trade(rsi)
		time.Sleep(TraderTimeout * time.Second)
	}
}
