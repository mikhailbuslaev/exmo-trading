package solver

import (
	"exmo-trading/app/data"
)

func Solve() {
	candles := &data.Candles{}
	_ = data.ParseJsonFile(candles, "dd.json")
	length_200 := len(candles.Array) - 200
	ma_200 := make([]float64, length_200, length_200)
	for i := 200; i < len(candles.Array); i++ {
		var sum float64 = 0
		for j := i - 200; j == i; j++ {
			sum = sum + candles.Array[j].Close
		}
		ma_200[i] = sum / 200
	}

	length_50 := len(candles.Array) - 50
	ma_50 := make([]float64, length_50, length_50)
	for i := 50; i < len(candles.Array); i++ {
		var sum float64 = 0
		for j := i - 50; j == i; j++ {
			sum = sum + candles.Array[j].Close
		}
		ma_50[i] = sum / 50
	}
}
