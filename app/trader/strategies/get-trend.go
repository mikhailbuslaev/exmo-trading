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

func (t *GetTrend) TestStrategy(c *data.Candles, ma []float64) {
	length := len(ma)
	lastOrderPrice := c.Array[0].Close
	profit := 0.00
	sumProfit := 0.00
	for i := 5; i < length; i++ {
		if c.Array[i].Close > ma[i] && c.Array[i-1].Close < ma[i-1]{
			profit = lastOrderPrice - c.Array[i].Close
			sumProfit = profit + sumProfit - 0.1*0.003*(lastOrderPrice+c.Array[i].Close)
			fmt.Println("profit is " + fmt.Sprintf("%f", profit))
			lastOrderPrice = c.Array[i].Close
		}
		if c.Array[i].Close < ma[i] && c.Array[i-1].Close > ma[i-1] {
			profit = -lastOrderPrice + c.Array[i].Close
			sumProfit = profit + sumProfit - 0.1*0.003*(lastOrderPrice+c.Array[i].Close)
			fmt.Println("profit is " + fmt.Sprintf("%f", profit))
			lastOrderPrice = c.Array[i].Close
		}
	}
	fmt.Println(sumProfit)
}

func (t *GetTrend) TestAnalyze() error {
	candles := &data.Candles{}
	candles.Array = make([]data.Candle, 0, 2000)
	err := candles.Read(t.CandlesFile)
	if err != nil {
		return err
	}

	priceArray := traderutils.GetArrayFromCandles(candles)
	ma := traderutils.GetMA(priceArray, t.MAFrame)
	t.TestStrategy(candles, ma)
	return nil
}
