package strategies

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/traderutils"
	"fmt"
)

type GetTrend struct {
	CandlesFile string
	MAFrame     int
}

func (t *GetTrend) Set(candlesFile string) {
	t.CandlesFile = candlesFile
	t.MAFrame = 200
}

func (t *GetTrend) Solve(c *data.Candles, ma []float64) string {
	length := len(ma)
	if c.Array[length].Close > ma[length] {
		return "bull"
	}
	if c.Array[length].Close < ma[length] {
		return "bear"
	}
	return "empty"
}

func (t *GetTrend) Analyze() (string, error) {
	candles := &data.Candles{}
	candles.Array = make([]data.Candle, 0, 250)
	err := candles.Read(t.CandlesFile)
	if err != nil {
		return "empty", err
	}

	priceArray := traderutils.GetArrayFromCandles(candles)
	ma := traderutils.GetMA(priceArray, t.MAFrame)
	return t.Solve(candles, ma), nil
}
