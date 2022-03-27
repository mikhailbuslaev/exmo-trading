package candles

import (
	"encoding/json"
	"fmt"
	"io"
	"mikhailbuslaev/exmo/app/utils"
	"net/http"
)

func Load(pair string, from, to int64) {
	resp, err := http.Get("https://api.exmo.me/v1.1/candles_history?symbol=" + pair +
		"&resolution=30&from=" + fmt.Sprintf("%d", from) + "&to=" + fmt.Sprintf("%d", to))
	if err != nil {
		fmt.Println("Unable to load last kandles")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to load last kandles")
	}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		fmt.Println("Unable to load last kandles")
	}

	candles := data["candles"]
	out, err := json.Marshal(candles)
	if err != nil {
		fmt.Println("Unable to load last kandles")
	}
	utils.Record(out, "candles/candles.txt")
}
