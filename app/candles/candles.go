package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Candles struct {
	Array []Candle
}

type Candle struct {
	Time   int64
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume float64
}

type Data interface {
	Marshal() ([]byte, error)
	Unmarshal(buf []byte) error
	Append(candles *Candles)
}

func (c *Candle) Unmarshal(buf []byte) error {
	err := json.Unmarshal(buf, c)
	if err != nil {
		fmt.Println("Unable to unmarshal json with candle")
	}
	return err
}

func (c *Candle) Marshal() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Fail while converting struct to json")
	}
	return b, err
}

func (c *Candle) Append(candles *Candles) {
	candles.Array = append(candles.Array, *c)
}

func (c *Candles) Unmarshal(buf []byte) error {
	err := json.Unmarshal(buf, &c.Array)
	if err != nil {
		fmt.Println("Unable to unmarshal json with candles")
	}
	return err
}

func (c *Candles) Marshal() ([]byte, error) {
	b, err := json.Marshal(&c.Array)
	if err != nil {
		fmt.Println("Fail while converting struct to json")
	}
	return b, err
}

func (c *Candles) Append(candles *Candles) {
	for i := range c.Array {
		candles.Array = append(candles.Array, c.Array[i])
	}
}

func ParseJsonFile(d Data, fileName string) error {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Unable to read json file")
	}
	err = d.Unmarshal(buf)
	if err != nil {
		fmt.Println("Unable to unmarshal json file")
	}
	return err
}

func RewriteJsonFile(d Data, fileName string) error {

	err := os.Truncate(fileName, 0)
	if err != nil {
		fmt.Println("Unable to clear json file")
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Unable to open json file")
	}

	data, err := d.Marshal()

	_, err1 := f.Write(data)
	if err1 != nil {
		fmt.Println("Unable to write data into json file")
	}

	f.Close()
	return err
}

func AppendJsonFile(d Data, fileName string) error {
	candles := &Candles{}
	err := ParseJsonFile(d, fileName)
	d.Append(candles)
	err = RewriteJsonFile(candles, fileName)
	return err
}
