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

func (m *MAintersectionTrader) Prepare(c *data.Candles, mashort, malong []float64) *data.Order {

	length := len(mashort)
	for i := m.LongMAFrame + 1; i < length; i++ {
		if mashort[i] > malong[i] && malong[i-1] > mashort[i-1] {
			if traderutils.CheckActuality(c.Array[i].Time / 1000) {
				order := traderutils.MakeOrder("buy", c.Array[i].Close, c.Array[i].Time/1000)
				return order
			}
		}
		if mashort[i] < malong[i] && malong[i-1] < mashort[i-1] {
			if traderutils.CheckActuality(c.Array[i].Time / 1000) {
				order := traderutils.MakeOrder("sell", c.Array[i].Close, c.Array[i].Time/1000)
				return order
			}
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
	return m.Prepare(candles, mashort, malong), nil
}
