package main

import (
	"mikhailbuslaev/exmo/solver/candles"
	"time"
)

func main() {
	now := time.Now()
	candles.Load("BTC_USD", 1648258200, now.Unix())
}
