package strategies

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/signals"
	"exmo-trading/app/trader/traderutils"
)

type GetTrend struct { // get trend gives bull or bear signals when actual price lower or hihger than ma200
	CandlesFile       string
	MAFrame           int
	CandlesFileVolume int
}

func (t *GetTrend) Set(candlesFile string, candlesFileVolume int) {
	t.CandlesFileVolume = candlesFileVolume
	t.CandlesFile = candlesFile
	t.MAFrame = 200
}

func (t *GetTrend) Solve(c *data.Candles, ma []float64) string {
	length := len(ma)
	if length != 0 {
		if c.Array[length].Close > ma[length] {
			return signals.Bull
		}
		if c.Array[length].Close < ma[length] {
			return signals.Bear
		}
	}
	return signals.NoSignals
}

func (t *GetTrend) Analyze() (string, error) {
	candles := &data.Candles{}
	candles.Array = make([]data.Candle, 0, 250)
	err := candles.Read(t.CandlesFile)
	if err != nil {
		return signals.NoSignals, err
	}

	priceArray := traderutils.GetArrayFromCandles(candles)
	ma := traderutils.GetMA(priceArray, t.MAFrame)
	return t.Solve(candles, ma), nil
}
