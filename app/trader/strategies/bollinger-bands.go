package strategies

import (
	"exmo-trading/app/data"
	"exmo-trading/app/trader/signals"
	"exmo-trading/app/trader/traderutils"
)

type BollingerBandsTrader struct { // rsi strategy gives long or short signals when rsi index goes lower than 30 or higher than 70
	CandlesFile       string
	CandlesFileVolume int
	Period            int
	Factor			  int
}

func (bb *BollingerBandsTrader) Set(candlesFile string, candlesFileVolume int) {
	bb.Period = 40
	bb.Factor = 2
	bb.CandlesFile = candlesFile
	bb.CandlesFileVolume = candlesFileVolume
}

func (bb *BollingerBandsTrader) Solve(c *data.Candles, topborder, bottomborder []float64) string {
	length := len(c.Array)
	if c.Array[length-1].Close > topborder[length-1] {
		return signals.Short
	}
	if c.Array[length-1].Close < bottomborder[length-1] {
		return signals.Long
	}
	return signals.NoSignals
}

func (bb *BollingerBandsTrader) Analyze() (string, error) {
	candles := &data.Candles{}
	candles.Array = make([]data.Candle, 0, bb.CandlesFileVolume)
	err := candles.Read(bb.CandlesFile)
	if err != nil {
		return signals.NoSignals, err
	}
	priceArray := traderutils.GetArrayFromCandles(candles)
	ma := traderutils.GetMA(priceArray, bb.Period)
	sd := traderutils.GetSD(priceArray, ma, bb.Period)
	length := len(ma)
	topborder := make([]float64, length)
	bottomborder := make([]float64, length)
	for i := 0; i < length; i++ {
		topborder[i] = ma[i] + sd[i]*float64(bb.Factor)
		bottomborder[i] = ma[i] - sd[i]*float64(bb.Factor)
	}
	return bb.Solve(candles, topborder, bottomborder), nil
}
