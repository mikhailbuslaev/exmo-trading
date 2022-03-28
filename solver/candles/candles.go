package candles

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mikhailbuslaev/exmo/app/utils"
	"net/http"
)

type Candle struct {
	Time   int64   `json:"t"`
	Open   float64 `json:"o"`
	Close  float64 `json:"c"`
	High    float64 `json:"h"`
	Low    float64 `json:"l"`
	Volume float64 `json:"v"`
}

func Load(pair string, from, to int64) {
	resp, err := http.Get("https://api.exmo.me/v1.1/candles_history?symbol=" + pair +
		"&resolution=5&from=" + fmt.Sprintf("%d", from) + "&to=" + fmt.Sprintf("%d", to))

	if err != nil {
		fmt.Println("Unable to load last candles")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to read last candles")
	}
	var data map[string][]Candle
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println("Unable to unmarshal last candles")
	}

	candles := data["candles"]

	for i := range candles {
		out, err := json.Marshal(candles[i])
		if err != nil {
			fmt.Println("Unable to marshal candles")
		}
		if i == 0 {
			utils.Record([]byte(","), "candles/candles.json")
			utils.RecordNewLine("candles/candles.json")
		}
		if i == len(candles) - 1 {
			utils.Rewrite(out, "candles/last-candle.json")
			utils.Record(out, "candles/candles.json")
		} else {
			utils.Record(out, "candles/candles.json")
			utils.Record([]byte(","), "candles/candles.json")
			utils.RecordNewLine("candles/candles.json")
		}

		
		fmt.Println("Candle loaded")

	}
}

func FindLastCandle() int64 {

	buf, err := ioutil.ReadFile("candles/last-candle.json")
	candle := Candle{}
	err = json.Unmarshal(buf, &candle)
	if err != nil {
		fmt.Println("Unable to unmarshal line from last-candle.json")
	}
	
	return candle.Time
}
