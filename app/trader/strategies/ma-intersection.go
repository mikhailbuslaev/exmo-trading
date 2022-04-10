package strategies

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/traderutils"
)

type MAintersectionTrader struct {
	CandlesFile  string
	LongMAFrame  int
	ShortMAFrame int
}

func (m *MAintersectionTrader) Set() {
	m.CandlesFile = "cache/5min-candles.csv"
	m.LongMAFrame = 200
	m.ShortMAFrame = 50
}

func (m *MAintersectionTrader) Solve(c *data.Candles, mashort, malong []float64) *data.Order {
	length := len(mashort)
	if mashort[length] > malong[length] && malong[length-1] > mashort[length-1] {
		if traderutils.CheckActuality(c.Array[length].Time / 1000) {
			order := traderutils.MakeOrder("long", c.Array[length].Close, c.Array[length].Time/1000)
			return order
		}
	}
	if mashort[length] < malong[length] && malong[length-1] < mashort[length-1] {
		if traderutils.CheckActuality(c.Array[length].Time / 1000) {
			order := traderutils.MakeOrder("long", c.Array[length].Close, c.Array[length].Time/1000)
			return order
		}
	}
	return nil
}

func (m *MAintersectionTrader) Analyze() (*data.Order, error) {
	candles := &data.Candles{}
	candles.Array = make([]data.Candle, 0, 250)
	err := candles.Read(m.CandlesFile)

	if err != nil {
		return nil, err
	}

	priceArray := traderutils.GetArrayFromCandles(candles)
	malong := traderutils.GetMA(priceArray, m.LongMAFrame)
	mashort := traderutils.GetMA(priceArray, m.ShortMAFrame)
	return m.Solve(candles, mashort, malong), nil
}
