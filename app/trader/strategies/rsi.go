package strategies

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/traderutils"
)

type RSItrader struct {
	CandlesFile string
	Period      int
}

func (rsi *RSItrader) Set() {
	rsi.CandlesFile = "cache/5min-candles.csv"
	rsi.Period = 14
}

func (rsi *RSItrader) Solve(c *data.Candles, avggain, avglose []float64) *data.Order {
	length := len(avggain)
	index := 100 - (100 / (1 + (avggain[length] / avglose[length])))
	if index > 70 {
		order := traderutils.MakeOrder("short", c.Array[length].Close, c.Array[length].Time/1000)
		return order
	}
	if index < 30 {
		order := traderutils.MakeOrder("long", c.Array[length].Close, c.Array[length].Time/1000)
		return order
	}
	return nil
}

func (rsi *RSItrader) Analyze() (*data.Order, error) {
	candles := &data.Candles{}
	candles.Array = make([]data.Candle, 0, 250)
	err := candles.Read(rsi.CandlesFile)
	if err != nil {
		return nil, err
	}
	priceArray := traderutils.GetArrayFromCandles(candles)
	avggain, avglose := traderutils.CountAvgChanges(priceArray, rsi.Period)
	return rsi.Solve(candles, avggain, avglose), nil
}
