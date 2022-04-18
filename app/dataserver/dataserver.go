package dataserver //dataserver is autonomic microservice loading and updating candles every 60 seconds
// at this stage of project dataserver working with only 5-min candles, but amount of files can be expanded

import (
	"exmo-trading/app/data"
	"exmo-trading/app/query"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Handler struct {
	Symbol        		string 			`yaml:"Symbol"`
	Resolution    		int64  			`yaml:"Resolution"`
	CandlesFile   		string 			`yaml:"CandlesFile"`
	CandlesVolume 		int64  			`yaml:"CandlesVolume"`
	DataServerTimeout 	time.Duration 	`yaml:"DataServerTimeout"`
}

func (h *Handler) Nothing() {
	
}

func (h *Handler) LoadCandles(from, to string) error {
	stringResolution := fmt.Sprintf("%d", h.Resolution)
	q := query.GetQuery{Method: "candles_history?symbol=" + h.Symbol +
		"&resolution=" + stringResolution + "&from=" + from + "&to=" + to}

	resp, err := query.Exec(&q)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	candles := &data.Candles{}
	candles.Array = make([]data.Candle, 0, h.CandlesVolume)

	err = candles.ParseJson([]byte(body))
	if err != nil {
		return err
	}

	if len(candles.Array) > 0 {
		err = candles.Write(h.CandlesFile)
		if err != nil {
			return err
		}
		fmt.Println("dataserver: new candles loaded")
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

func (h *Handler) InitCandles() error {
	err := ClearFile(h.CandlesFile)
	if err != nil {
		return err
	}
	t := time.Now()

	err = h.LoadCandles(fmt.Sprintf("%d", t.Unix()-60*h.Resolution*h.CandlesVolume-1), fmt.Sprintf("%d", t.Unix()+1))
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Run() {
	for {
		err := h.InitCandles()
		time.Sleep(h.DataServerTimeout * time.Second)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
