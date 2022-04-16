package data

import (
	"database/sql"
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

type Orders struct {
	Array []Order `json:"orders"`
}

type Order struct {
	Id     int64   `json:"i"`
	Action string  `json:"a"`
	Time   int64   `json:"t"`
	Price  float64 `json:"p"`
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
	err := json.Unmarshal(buf, c)
	if err != nil {
		fmt.Println("json parsing fail")
		return err
	}
	return nil
}

func (o *Order) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, o)
	if err != nil {
		return err
	}
	return nil
}

func (o *Orders) ParseJson(buf []byte) error {
	err := json.Unmarshal(buf, o)
	if err != nil {
		fmt.Println("json parsing fail")
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

func (o *Order) ParseString(input []string) error {
	var err error
	o.Time, err = strconv.ParseInt(input[0], 10, 64)
	if err != nil {
		return err
	}
	o.Action = input[1]
	o.Time, err = strconv.ParseInt(input[3], 10, 64)
	if err != nil {
		return err
	}
	o.Price, err = strconv.ParseFloat(input[3], 64)
	if err != nil {
		return err
	}
	return nil
}

func (o *Order) MakeString() []string {
	output := make([]string, 4)
	output[0] = fmt.Sprintf("%d", o.Id)
	output[1] = o.Action
	output[2] = fmt.Sprintf("%f", o.Price)
	output[3] = fmt.Sprintf("%d", o.Time)
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

func (o *Orders) Read(fileName string) error {

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
		order := &Order{}
		err := order.ParseString(records[i])
		if err != nil {
			return err
		}
		o.Array = append(o.Array, *order)
	}
	return nil
}

func (o *Order) Read(fileName string) error {

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
	err = o.ParseString(record)
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

func (o *Orders) Write(fileName string) error {

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	w := csv.NewWriter(f)

	defer w.Flush()

	for i := range o.Array {
		err := w.Write(o.Array[i].MakeString())
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Order) Write(fileName string) error {

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	w := csv.NewWriter(f)

	defer w.Flush()

	err = w.Write(o.MakeString())
	if err != nil {
		return err
	}
	return nil
}

func (o *Order) Print() {
	fmt.Println(o.Action + ", " + fmt.Sprintf("%d", o.Time) + ", " + fmt.Sprintf("%f", o.Price))
}

func (c *Candles) ScanRows(rows *sql.Rows) error {
	i := 0
	for rows.Next() {
		err := rows.Scan(&c.Array[i].Time, &c.Array[i].Open,
			&c.Array[i].Close, &c.Array[i].High, &c.Array[i].Low, &c.Array[i].Volume)
		if err != nil {
			return err
		}
		i++
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
