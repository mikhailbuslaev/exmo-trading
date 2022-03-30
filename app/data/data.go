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

type Decisions struct {
	Array []Decision
}

type Decision struct {
	Time    int64
	Thought string
}

type Data interface {
	Marshal() ([]byte, error)
	Unmarshal(buf []byte) error
	AppendToCandles(c *Candles)
	AppendToDecisions(d *Decisions)
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

func (c *Candle) AppendToCandles(candles *Candles) {
	candles.Array = append(candles.Array, *c)
}

func (d *Candle) AppendToDecisions(decisions *Decisions) {
	fmt.Println("Forbidden method")
}

func (c *Candles) Unmarshal(buf []byte) error {
	err := json.Unmarshal(buf, &c.Array)
	if err != nil {
		fmt.Println("Unable to unmarshal json with candles")
	}
	return err
}

func (c *Candles) Marshal() ([]byte, error) {
	b, err := json.Marshal(c.Array)
	if err != nil {
		fmt.Println("Fail while converting struct to json")
	}
	return b, err
}

func (c *Candles) AppendToCandles(candles *Candles) {
	for i := range c.Array {
		candles.Array = append(candles.Array, c.Array[i])
	}
}

func (c *Candles) AppendToDecisions(decisions *Decisions) {
	fmt.Println("Forbidden method")
}

func (d *Decision) Unmarshal(buf []byte) error {
	err := json.Unmarshal(buf, d)
	if err != nil {
		fmt.Println("Unable to unmarshal json with decision")
	}
	return err
}

func (d *Decision) Marshal() ([]byte, error) {
	b, err := json.Marshal(d)
	if err != nil {
		fmt.Println("Fail while converting struct to json")
	}
	return b, err
}

func (d *Decision) AppendToDecisions(decisions *Decisions) {
	decisions.Array = append(decisions.Array, *d)
}

func (d *Decision) AppendToCandles(decisions *Decisions) {
	fmt.Println("Forbidden method")
}

func (d *Decisions) Unmarshal(buf []byte) error {
	err := json.Unmarshal(buf, &d.Array)
	if err != nil {
		fmt.Println("Unable to unmarshal json with decision")
	}
	return err
}

func (d *Decisions) Marshal() ([]byte, error) {
	b, err := json.Marshal(d.Array)
	if err != nil {
		fmt.Println("Fail while converting struct to json")
	}
	return b, err
}

func (d *Decisions) AppendToDecisions(decisions *Decisions) {
	for i := range d.Array {
		decisions.Array = append(decisions.Array, d.Array[i])
	}
}

func (d *Decisions) AppendToCandles(candles *Candles) {
	fmt.Println("Forbidden method")
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

func AppendJsonFile(d Data, fileName string, flag string) error {
	if flag == "candles" {
		output := &Candles{}
		err := ParseJsonFile(d, fileName)
		d.AppendToCandles(output)
		err = RewriteJsonFile(output, fileName)
		return err
	} else {
		output := &Decisions{}
		err := ParseJsonFile(output, fileName)
		d.AppendToDecisions(output)
		err = RewriteJsonFile(output, fileName)
		return err
	}
}
