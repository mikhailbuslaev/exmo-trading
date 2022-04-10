package traderutils

import (
	"exmo-trading/app/data"
	"time"
)

func GetMA(array []float64, frame int) []float64 {
	length := len(array)
	ma := make([]float64, length)
	for i := frame; i < length; i++ {
		sum := 0.00
		for j := i - frame; j < i; j++ {
			sum = sum + array[j]
		}
		ma[i] = sum / float64(frame)
	}
	return ma
}

func CheckActuality(input int64) bool {
	t := time.Now()
	return t.Unix()-input < 59
}

func GetArrayFromCandles(c *data.Candles) []float64 {
	length := len(c.Array)
	array := make([]float64, length)
	for i := range array {
		array[i] = c.Array[i].Close
	}
	return array
}

func ConvertCandleTime(inputTime int64) time.Time {
	t := inputTime / 1000
	return time.Unix(t, 0)
}

func MakeOrder(action string, price float64, time int64) *data.Order {
	o := &data.Order{}
	o.Action = action
	o.Price = price
	o.Time = time
	return o
}
