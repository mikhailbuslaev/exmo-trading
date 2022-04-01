package data

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Candles struct {
	Array []Candle `json:"candles"`
}

type Candle struct {
	Time   int64   `json:"t"`
	Open   float64 `json:"o"`
	Close  float64 `json:"c"`
	High   float64 `json:"h"`
	Low    float64 `json:"l"`
	Volume float64 `json:"v"`
}

type Data interface {
	ParseJson([]byte) error
	Read(string) error
	Write(string) error
}

func (c *Candle) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Candles) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, &c.Array)
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
		err := c.Array[i].ParseString(records[i])
		if err != nil {
			return err
		}
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
