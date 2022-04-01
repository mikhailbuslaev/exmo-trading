package main

import (
	"exmo-trading/app/data"
	"exmo-trading/app/query"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func LoadCandles(resolution, from, to, candlesFile, lastCandleFile string) error {
	q := query.GetQuery{Method: "candles_history?symbol=BTC_USD" +
		"&resolution=" + resolution + "&from=" + from + "&to=" + to}
	resp, err := query.Exec(&q)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	candles := &data.Candles{}

	err = candles.ParseJson([]byte(body))
	if err != nil {
		return err
	}

	err = candles.Write(candlesFile)
	if err != nil {
		return err
	}

	err = data.Rewrite(&candles.Array[len(candles.Array)], lastCandleFile)
	if err != nil {
		return err
	}

	return nil
}

func ClearFile(fileName string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCandles(resolution, candlesFile, lastCandleFile string) error {
	candle := &data.Candle{}
	err := candle.Read(lastCandleFile)
	if err != nil {
		return err
	}
	t := time.Now()
	err = LoadCandles(resolution, fmt.Sprintf("%d", candle.Time+1),
		fmt.Sprintf("%d", t.Unix()), candlesFile, lastCandleFile)
	if err != nil {
		return err
	}
	return nil
}

func InitCandles(resolution, candlesFile, lastCandleFile string) error {
	err := ClearFile(candlesFile)
	if err != nil {
		return err
	}
	t := time.Now()
	intresol, err := strconv.ParseInt(resolution, 10, 64)
	if err != nil {
		return err
	}
	err = LoadCandles(resolution, fmt.Sprintf("%d", t.Unix()),
		fmt.Sprintf("%d", t.Unix()-60*intresol*1000), candlesFile, lastCandleFile)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	for {

		time.Sleep(30 * time.Second)
	}
}
