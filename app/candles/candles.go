package candles

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
	ParseJsonFile(fileName string) error
	ParseJson(buf []byte) error
	AppendJsonFile(fileName string) error
	RewriteJsonFile(fileName string) error
	MakeJson() ([]byte, error)
}

func (c *Candle) MakeJson() ([]byte, error) {
	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println("Fail while convering struct to json")
	}
	return b, err
}

func (c *Candles) MakeJson() ([]byte, error) {
	b, err := json.Marshal(&c.Array)
	if err != nil {
		fmt.Println("Fail while convering struct to json")
	}
	return b, err
}

func (c *Candle) ParseJsonFile(fileName string) error {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Unable to read json file with candle")
	}
	err = json.Unmarshal(buf, c)
	if err != nil {
		fmt.Println("Unable to unmarshal json file with candle")
	}
	return err
}

func (c *Candles) ParseJsonFile(fileName string) error {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Unable to read json file with candles")
	}
	err = json.Unmarshal(buf, &c.Array)
	if err != nil {
		fmt.Println("Unable to unmarshal json file with candles")
	}
	return err
}

func (c *Candle) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, c)
	if err != nil {
		fmt.Println("Unable to unmarshal json with candle")
	}
	return err
}

func (c *Candles) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, &c.Array)
	if err != nil {
		fmt.Println("Unable to unmarshal json with candles")
	}
	return err
}

func (c *Candles) RewriteJsonFile(fileName string) error {

	err := os.Truncate(fileName, 0)
	if err != nil {
		fmt.Println("Unable to clear json file")
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Unable to open json file")
	}

	data, err := c.MakeJson()

	_, err1 := f.Write(data)
	if err1 != nil {
		fmt.Println("Unable to write data into json file")
	}

	f.Close()
	return err
}

func (c *Candle) RewriteJsonFile(fileName string) error {

	err := os.Truncate(fileName, 0)
	if err != nil {
		fmt.Println("Unable to clear json file")
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Unable to open json file")
	}

	data, err := c.MakeJson()

	_, err1 := f.Write(data)
	if err1 != nil {
		fmt.Println("Unable to write data into json file")
	}

	f.Close()
	return err
}

func (c *Candle) AppendJsonFile(fileName string) error {
	candles := &Candles{}
	err := candles.ParseJsonFile(fileName)
	candles.Array = append(candles.Array, *c)
	err = candles.RewriteJsonFile(fileName)
	return err
}

func (c *Candles) AppendJsonFile(fileName string) error {
	candles := &Candles{}
	err := candles.ParseJsonFile(fileName)
	for i := range c.Array {
		candles.Array = append(candles.Array, c.Array[i])
	}
	err = candles.RewriteJsonFile(fileName)
	return err
}

func AppendToFile(d Data, fileName string) {
	fmt.Println(d.AppendJsonFile(fileName))
}

func RewriteFile(d Data, fileName string) {
	fmt.Println(d.RewriteJsonFile(fileName))
}

func ParseFile(d Data, fileName string) {
	fmt.Println(d.ParseJsonFile(fileName))
}

func ParseJson(d Data, buf []byte) {
	fmt.Println(d.ParseJson(buf))
}

func MakeJson(d Data) ([]byte, error) {
	return d.MakeJson()
}
