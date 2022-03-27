package main

import (
	"fmt"
	"mikhailbuslaev/exmo/solver/candles"
	"time"
)

func main() {
	for 1 == 1 {
		now := time.Now()
		from := candles.FindLastCandle()
		if from == 0 {
			from = now.Unix() - 604800
		}
		fmt.Println(from)
		candles.Load("BTC_USD", from, now.Unix())
		time.Sleep(60 * time.Second)
	}
}
