package strategies

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/signals"
	"exmo-trading/app/trader/traderutils"
	"fmt"
)

type RSItrader struct { // rsi strategy gives long or short signals when rsi index goes lower than 30 or higher than 70
	CandlesFile string
	Period      int
}

func (rsi *RSItrader) Set(candlesFile string) {
	rsi.CandlesFile = candlesFile
	rsi.Period = 14
}

func (rsi *RSItrader) Solve(c *data.Candles, avggain, avglose []float64) string {
	length := len(avggain)
	index := 100 - (100 / (1 + (avggain[length-1] / avglose[length-1])))
	fmt.Println("RSI is " + fmt.Sprintf("%f", index))
	if index > 70 {
		return signals.Short
	}
	if index < 30 {
		return signals.Long
	}
	return signals.NoSignals
}

func (rsi *RSItrader) Analyze() (string, error) {
	candles := &data.Candles{}
	candles.Array = make([]data.Candle, 0, 250)
	err := candles.Read(rsi.CandlesFile)
	if err != nil {
		return signals.NoSignals, err
	}
	priceArray := traderutils.GetArrayFromCandles(candles)
	avggain, avglose := traderutils.CountAvgChanges(priceArray, rsi.Period)
	return rsi.Solve(candles, avggain, avglose), nil
}
