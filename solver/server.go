package main

import (
	"fmt"
	"mikhailbuslaev/exmo/solver/candles"
	"time"
)

func main() {
	for 1 == 1 {
		now := time.Now()
		from := candles.FindLastCandle()/1000
		if from == 0 {
			from = now.Unix() - 86400
		} else {
			from = from + 1
		}
		candles.Load("BTC_USD", from, now.Unix())
		fmt.Println("Waiting...")
		time.Sleep(30 * time.Second)
	}
}
