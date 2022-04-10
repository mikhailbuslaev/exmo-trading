package main

import (
	"exmo-trading/app/dataserver"
	"exmo-trading/app/trader"
)

func main() {
	dataserver.Launch()
	trader.Launch()

}
