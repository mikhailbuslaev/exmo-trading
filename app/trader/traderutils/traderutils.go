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
		for j := i - frame; j <= i; j++ {
			sum = sum + array[j]
		}
		ma[i] = sum / float64(frame)
	}
	return ma
}

func CountAvgChanges(array []float64, frame int) ([]float64, []float64) {
	length := len(array)
	avggain := make([]float64, 0, length)
	avglose := make([]float64, 0, length)

	for i := frame; i < length; i++ {
		n, m := 0, 0
		avggain[i], avglose[i] = 0.00, 0.00
		loses, gains := 0.00, 0.00
		for j := i - frame; j <= i; j++ {
			if array[j] < array[j-1] {
				loses = loses + (array[j-1]-array[j])/array[j-1]
				n++
			}
			if array[i] > array[i-1] {
				gains = gains + (array[j]-array[j-1])/array[j-1]
				m++
			}
		}

		if m != 0 {
			avggain[i] = gains / float64(m)
		}

		if n != 0 {
			avglose[i] = loses / float64(n)
		}
	}
	return avggain, avglose
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
