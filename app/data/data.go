package data // this package define and descibes behaviour of main data types in app

import (
	"encoding/csv"  // all data stored in csv cache ("cache/file.csv")
	"encoding/json" // need for parsing candles and responses
	"fmt"           // printing errors in console
	"os"            // working with cache files(reading, writing, rewriting)
	"strconv"       // need for working with csv
)

type Candles struct {
	Array []Candle `json:"candles"` // candles array
}

type Candle struct {
	Time   int64   `json:"t"` // close time of candle
	Open   float64 `json:"o"` // open price of candle
	Close  float64 `json:"c"` // close price of candle
	High   float64 `json:"h"` // high price of candle
	Low    float64 `json:"l"` // low price of candle
	Volume float64 `json:"v"` // traded volume of candle
}

type Trades struct {
	Array []Trade `json:"trades"` // trades array
}

type Trade struct {
	Id         int64   `json:"id"`          // id of trade
	Action     string  `json:"action"`      // long or short
	OpenPrice  float64 `json:"open-price"`  // starting price of trade
	ClosePrice float64 `json:"close-price"` // finished price of trade
	Quantity   float64 `json:"quantity"`    // trade volume, example : amount of usdt when you buy btc
	StopLimit  float64 `json:"stop-limit"`  // when price goes to stop limit, trade closed
	Status     string  `json:"status"`      // opened or closed
}

type Ticker struct {
	BuyPrice  string `json:"buy_price"` // ticker need for getting actual price of trading pair
	SellPrice string `json:"sell_price"`
	LastTrade string `json:"last_trade"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Avg       string `json:"avg"`
	Vol       string `json:"vol"`
	VolCurr   string `json:"vol_curr"`
	Updated   int64  `json:"updated"`
}

type Data interface {
	ParseJson([]byte) error
	Read(string) error
	Write(string) error
}

func (t *Ticker) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Ticker) ParseJsonTickers(buf []byte, pair string) error {
	var tickers map[string]json.RawMessage
	err := json.Unmarshal(buf, &tickers)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tickers[pair], t)
	if err != nil {
		return err
	}
	return nil
}

func (c *Candle) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Candles) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trade) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trades) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, t)
	if err != nil {
		return err
	}
	return nil
}

func (c *Candle) ParseString(input []string) error {
	var err error
	c.Time, err = strconv.ParseInt(input[0], 10, 64)
	if err != nil {
		return err
	}
	c.Open, err = strconv.ParseFloat(input[1], 64)
	if err != nil {
		return err
	}
	c.Close, err = strconv.ParseFloat(input[2], 64)
	if err != nil {
		return err
	}
	c.High, err = strconv.ParseFloat(input[3], 64)
	if err != nil {
		return err
	}
	c.Low, err = strconv.ParseFloat(input[4], 64)
	if err != nil {
		return err
	}
	c.Volume, err = strconv.ParseFloat(input[5], 64)
	if err != nil {
		return err
	}
	return nil
}

func (t *Trade) ParseString(input []string) error {
	var err error
	t.Id, err = strconv.ParseInt(input[0], 10, 64)
	if err != nil {
		return err
	}
	t.Action = input[1]
	t.Quantity, err = strconv.ParseFloat(input[2], 64)
	if err != nil {
		return err
	}
	t.OpenPrice, err = strconv.ParseFloat(input[3], 64)
	if err != nil {
		return err
	}
	t.ClosePrice, err = strconv.ParseFloat(input[4], 64)
	if err != nil {
		return err
	}
	t.StopLimit, err = strconv.ParseFloat(input[5], 64)
	if err != nil {
		return err
	}
	t.Status = input[6]
	return nil
}

func (t *Trade) MakeString() []string {
	output := make([]string, 7)
	output[0] = fmt.Sprintf("%d", t.Id)
	output[1] = t.Action
	output[2] = fmt.Sprintf("%f", t.Quantity)
	output[3] = fmt.Sprintf("%f", t.OpenPrice)
	output[4] = fmt.Sprintf("%f", t.ClosePrice)
	output[5] = fmt.Sprintf("%f", t.StopLimit)
	output[6] = t.Status
	return output
}

func (c *Candle) MakeString() []string {
	output := make([]string, 6)
	output[0] = fmt.Sprintf("%d", c.Time)
	output[1] = fmt.Sprintf("%f", c.Open)
	output[2] = fmt.Sprintf("%f", c.Close)
	output[3] = fmt.Sprintf("%f", c.High)
	output[4] = fmt.Sprintf("%f", c.Low)
	output[5] = fmt.Sprintf("%f", c.Volume)
	return output
}

func (c *Candles) Read(fileName string) error {

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer f.Close()
	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		return err
	}
	for i := range records {
		candle := &Candle{}
		err := candle.ParseString(records[i])
		if err != nil {
			return err
		}
		c.Array = append(c.Array, *candle)
	}
	return nil
}

func (c *Candle) Read(fileName string) error {

	f, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()
	r := csv.NewReader(f)

	record, err := r.Read()
	if err != nil {
		return err
	}
	err = c.ParseString(record)
	if err != nil {
		return err
	}

	return nil
}

func (t *Trades) Read(fileName string) error {

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer f.Close()
	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		return err
	}
	for i := range records {
		trade := &Trade{}
		err := trade.ParseString(records[i])
		if err != nil {
			return err
		}
		t.Array = append(t.Array, *trade)
	}
	return nil
}

func (t *Trade) Read(fileName string) error {

	f, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()
	r := csv.NewReader(f)

	record, err := r.Read()
	if err != nil {
		return err
	}
	err = t.ParseString(record)
	if err != nil {
		return err
	}

	return nil
}

func (c *Candles) Write(fileName string) error {

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	w := csv.NewWriter(f)

	defer w.Flush()

	for i := range c.Array {
		err := w.Write(c.Array[i].MakeString())
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Candle) Write(fileName string) error {

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	w := csv.NewWriter(f)

	defer w.Flush()

	err = w.Write(c.MakeString())
	if err != nil {
		return err
	}
	return nil
}

func (t *Trades) Write(fileName string) error {

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	w := csv.NewWriter(f)

	defer w.Flush()

	for i := range t.Array {
		err := w.Write(t.Array[i].MakeString())
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Trade) Write(fileName string) error {

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	w := csv.NewWriter(f)

	defer w.Flush()

	err = w.Write(t.MakeString())
	if err != nil {
		return err
	}
	return nil
}

func Rewrite(d Data, fileName string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()
	err = f.Truncate(0)
	if err != nil {
		return err
	}

	err = d.Write(fileName)
	if err != nil {
		return err
	}
	return nil
}
