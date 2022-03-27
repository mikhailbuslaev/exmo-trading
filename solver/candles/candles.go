package candles

import (
	"encoding/json"
	"fmt"
	"io"
	"mikhailbuslaev/exmo/app/utils"
	"net/http"
	"os"
)

type Candle struct {
	Time   int64   `json:"t"`
	Open   float64 `json:"o"`
	Close  float64 `json:"c"`
	Max    float64 `json:"h"`
	Min    float64 `json:"l"`
	Volume float64 `json:"v"`
}

func Load(pair string, from, to int64) {
	resp, err := http.Get("https://api.exmo.me/v1.1/candles_history?symbol=" + pair +
		"&resolution=30&from=" + fmt.Sprintf("%d", from) + "&to=" + fmt.Sprintf("%d", to))
	if err != nil {
		fmt.Println("Unable to load last candles")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to load last candles")
	}
	var data map[string][]Candle
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println("Unable to load last candles")
	}

	candles := data["candles"]

	for i := range candles {
		out, err := json.Marshal(candles[i])
		if err != nil {
			fmt.Println("Unable to load last candles")
		}
		utils.Record(out, "candles/candles.txt")
		fmt.println("---")
	}
}

func FindLastCandle() int64 {

	file, err := os.Open("candles/candles.txt")
	if err != nil {
		fmt.Println("Unable to load last candles")
	}

	stat, err := os.Stat("candles/candles.txt")
	if err != nil {
		fmt.Println("Unable to load last candles")
	}

	start := stat.Size() - 86
	var buf []byte
	_, err = file.ReadAt(buf, start)
	if err != nil {
		fmt.Println("Unable to load last candles")
	}
	var candle Candle
	err = json.Unmarshal(buf, &candle)
	if err != nil {
		fmt.Println("Unable to load last candles")
	}
	return candle.Time
}
