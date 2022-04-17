package main //dataserver is autonomic microservice loading and updating candles every 60 seconds
// at this stage of project dataserver working with only 5-min candles, but amount of files can be expanded

import (
	"exmo-trading/app/data"
	"exmo-trading/app/query"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const CandlesFile string = "/cache/5min-candles.csv"
const DataServerTimeout time.Duration = 60

type Handler struct {
	Context Context `json:"context"`
}

type Context struct {
	Symbol        string `json:"symbol"`
	Resolution    int64  `json:"resolution"`
	CandlesFile   string `json:"candles-file"`
	DbTable       string `json:"db-table"`
	CandlesVolume int64  `json:"candles-volume"`
}

func (h *Handler) Set(candlesFile string) {
	h.Context.Symbol = "BTC_USDT"
	h.Context.Resolution = 5
	h.Context.CandlesFile = candlesFile
	h.Context.DbTable = "5min-candles"
	h.Context.CandlesVolume = 250
}

func (h *Handler) LoadCandles(from, to string) error {
	stringResolution := fmt.Sprintf("%d", h.Context.Resolution)
	q := query.GetQuery{Method: "candles_history?symbol=" + h.Context.Symbol +
		"&resolution=" + stringResolution + "&from=" + from + "&to=" + to}

	resp, err := query.Exec(&q)
	if err != nil {
		fmt.Println("no candles...")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error while reading response body")
		return err
	}

	defer resp.Body.Close()

	candles := &data.Candles{}
	candles.Array = make([]data.Candle, 0, h.Context.CandlesVolume)

	err = candles.ParseJson([]byte(body))
	if err != nil {
		fmt.Println("error while parsing json body")
		return err
	}

	if len(candles.Array) > 0 {
		err = candles.Write(h.Context.CandlesFile)
		if err != nil {
			fmt.Println("error while appending data")
			return err
		}
		fmt.Println("new candles loaded")
	} else {
		fmt.Println("no new candles")
	}

	return nil
}

func ClearFile(fileName string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("error while opening file")
		return err
	}

	defer f.Close()
	err = f.Truncate(0)
	if err != nil {
		fmt.Println("error while clearing file")
		return err
	}
	return nil
}

func (h *Handler) InitCandles() error {
	err := ClearFile(h.Context.CandlesFile)
	if err != nil {
		return err
	}
	t := time.Now()

	err = h.LoadCandles(fmt.Sprintf("%d", t.Unix()-60*h.Context.Resolution*h.Context.CandlesVolume), fmt.Sprintf("%d", t.Unix()))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	parent := filepath.Dir(filepath.Dir(wd))

	fivemin := &Handler{}
	fivemin.Set(parent + CandlesFile)
	for {
		err := fivemin.InitCandles()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("waiting new candles")
		time.Sleep(DataServerTimeout * time.Second)
	}
}
