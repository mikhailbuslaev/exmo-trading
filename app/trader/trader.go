package trader // trader is autonomic microservice that launch strategies and handling signals by making events

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

type Event struct { 		// event is need to making trades from strategy signals
							//we making trade relying on previous trades and last price
	Action  	string		//handled signal:long or short
	LastPrice	float64 	//actual price of traded pair
	Trades      data.Trades //array of open trades
}

type Trader struct {
	Analyzer
	Context TraderContext
}

type Analyzer interface {
	Analyze() (string, error) //strategies implement this interface
}

type TraderContext {
	TradingPair 	  string	`yaml:"TradingPair"`		//trading pair; exaple: "BTC_USDT"
	TradesFile        string	`yaml:"TradesFile"`			//trades file name; example: "/cache/trades.csv"
	TradesHistoryFile string	`yaml:"TradesHistoryFile"`	//trades history file name; example: "/cache/trades-history.csv"
	CandlesFile       string	`yaml:"CandlesFile"`		//candles file name; example: "/cache/5min-candles.csv"
	CandlesFileVolume int		`yaml:"CandlesFileVolume"`	// volume of candles in candles file; default: 250
	TradesFileVolume  int		`yaml:"TradesFileVolume"`	// max amount of open trades; default: 10 (means 10 trades)
	StopLimitPercent  float64 	`yaml:"StopLimitPercent"`	// 0.01 means 1%, when price change goes higher 
															// or lower this percent, we close unprofitable trade
															// default: 0.01
	TradeAmount 	  float64	`yaml:"TradeAmount"`		// trade volume; example: 100 USDT in BTC_USDT pair
	TraderTimeout time.Duration `yaml:"TraderTimeout"`		// time of trader reload; default: 30
}

func (t *TraderContext) Nothing() {

}

func (e *Event) GetLastPrice(ctx *TraderContext) error {
	q := &PostQuery{Method:"ticker"}
	resp, err := query.Exec(q)
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

	err = ticker.ParseJsonTickers([]byte(body), ctx.TradingPair)
	if err != nil {
		return 0, err
	}
	e.LastPrice, err = strconv.ParseFloat(ticker.Avg, 64)
	if err != nil {
		return 0, err
	}
	return nil
}

func (e *Event) Init(eventType string, ctx *TraderContext) error { // loading context for our event

	err := e.GetLastPrice(ctx)
	if err != nil {
		return err
	}

	trades := &data.Trades{}
	trades.Array = make([]data.Trade, 0, ctx.TradesFileVolume)
	err = trades.Read(ctx.TradesFile)
	if err != nil {
		return err
	}
	e.Trades = *trades
	e.Action = eventType
	return nil
}

func (e *Event) HandleOpenedTrades(ctx *TraderContext) error {
	for i := range e.Trades.Array {
		if e.Trades.Array[i].Action == signals.Long && e.Trades.Array[i].StopLimit > e.LastPrice {
			err := e.CloseTrade(&e.Trades.Array[i], ctx)
			if err != nil {
				return err
			}
		}
		if e.Trades.Array[i].Action == signals.Short && e.Trades.Array[i].StopLimit < e.LastPrice {
			err := e.CloseTrade(&e.Trades.Array[i], ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Event) HandleSignal(ctx *TraderContext) error {
	if e.Action == signals.Long || e.Action == signals.Short {
		if len(e.Trades.Array) != 0 && e.Trades.Array[len(e.Trades.Array)-1].Action != e.Action {
			lastTrade := &e.Trades.Array[len(e.Trades.Array)-1]
			if lastTrade.Action != e.Action {
				err := e.CloseTrade(lastTrade, ctx)
				if err != nil {
					return err
				}
				err = e.OpenTrade(ctx)
				if err != nil {
					return err
				}
			}
		}
		if len(e.Trades.Array) == 0 {
			err := e.OpenTrade(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *Event) Handle(ctx *TraderContext) error { // checking event context and making trades
	err := e.HandleOpenedTrades(ctx)
	if err != nil {
		return err
	}
	err = e.HandleSignal(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e *Event) OpenTrade(ctx *TraderContext) error {
	newTrade := &data.Trade{}
	newTrade.Action = e.Action
	now := time.Now()
	newTrade.Id = now.Unix()
	if e.Action == signals.Long {
		newTrade.StopLimit = e.Context.LastPrice - e.LastPrice*ctx.StopLimitPercent
	} else if e.Action == signals.Short {
		newTrade.StopLimit = e.Context.LastPrice + e.LastPrice*ctx.StopLimitPercent
	}
	newTrade.OpenPrice = e.LastPrice
	newTrade.Status = signals.TradeOpened
	newTrade.Quantity = ctx.TradeAmount
	e.Trades.Array = append(e.Trades.Array, *newTrade)
	err := data.Rewrite(&e.Trades, ctx.TradesFile)
	if err != nil {
		return err
	}
	return nil
}

func (e *Event) CloseTrade(t *data.Trade, ctx *TraderContext) error {
	var tradeIndex int
	for i := range e.Trades.Array {
		if t.Id == e.Trades.Array[i].Id {
			tradeIndex = i
		}
	}
	e.Trades.Array[tradeIndex].Status = signals.TradeClosed

	e.Trades.Array[tradeIndex].ClosePrice = e.LastPrice

	err = e.Trades.Array[tradeIndex].Write(ctx.TradesHistoryFile) // write trade to archive
	if err != nil {
		return err
	}

	e.Trades.Array = append(e.Trades.Array[:tradeIndex], e.Trades.Array[tradeIndex+1:]...) // deleting trade

	err = data.Rewrite(&e.Trades, ctx.TradesFile) // rewrite open trades file
	if err != nil {
		return err
	}
	return nil
}

func (t *Trader) Trade() error {
	eventType, err := t.Analyzer.Analyze()
	if err != nil {
		return err
	}

	e := &Event{}
	err = e.Init(eventType, &t.Context)
	if err != nil {
		return err
	}

	err = e.Handle(&t.Context)
	if err != nil {
		return err
	}
}

func (t *Trader) Run() {
	for {
		err := t.Trade()
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(t.Context.TraderTimeout * time.Second)
	}
}
