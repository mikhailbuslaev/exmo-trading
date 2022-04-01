package dataserver

import (
	"errors"
	"exmo-trading/app/data"
	"exmo-trading/app/database"
	"exmo-trading/app/query"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Context struct {
	Symbol         string
	Resolution     string
	CandlesFile    string
	LastCandleFile string
	DbTable        string
}

func (c *Context) LoadCandles(from, to string) error {
	q := query.GetQuery{Method: "candles_history?symbol=" + c.Symbol +
		"&resolution=" + c.Resolution + "&from=" + from + "&to=" + to}
	resp, err := query.Exec(&q)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	candles := &data.Candles{}

	err = candles.ParseJson([]byte(body))
	if err != nil {
		return err
	}

	err = candles.Write(c.CandlesFile)
	if err != nil {
		return err
	}

	err = data.Rewrite(&candles.Array[len(candles.Array)], c.LastCandleFile)
	if err != nil {
		return err
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

func (c *Context) UpdateCandles() error {
	candle := &data.Candle{}
	err := candle.Read(c.LastCandleFile)
	if err != nil {
		return err
	}
	t := time.Now()
	err = c.LoadCandles(fmt.Sprintf("%d", candle.Time+1), fmt.Sprintf("%d", t.Unix()))
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) InitCandles() error {
	err := ClearFile(c.CandlesFile)
	if err != nil {
		return err
	}
	t := time.Now()
	intresol, err := strconv.ParseInt(c.Resolution, 10, 64)
	if err != nil {
		return err
	}
	err = c.LoadCandles(fmt.Sprintf("%d", t.Unix()), fmt.Sprintf("%d", t.Unix()-60*intresol*1000))
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) SyncDbAndCache(dbConfigName string) error {
	db, err := database.Connect(dbConfigName)
	if err != nil {
		return err
	}

	rows, err := database.Select(db, "SELECT * FROM "+c.DbTable+" ORDER BY time ASC;")
	if err != nil {
		return err
	}

	cacheCandles := &data.Candles{}
	err = cacheCandles.Read(c.CandlesFile)
	if err != nil {
		return err
	}

	dbCandles := &data.Candles{}
	err = dbCandles.ScanRows(rows)
	if err != nil {
		return err
	}

	dbLen := len(dbCandles.Array)
	cacheLen := len(cacheCandles.Array)

	if dbLen == cacheLen {
		for i := 0; i < dbLen; i++ {
			if dbCandles.Array[i] != cacheCandles.Array[i] {
				dbCandles.Array[i] = cacheCandles.Array[i]
			}
		}
	} else {
		return errors.New("length of cache array not match length of db array")
	}

	err = database.Change(db, "DELETE * FROM "+c.DbTable+";")
	if err != nil {
		return err
	}

	for i := 0; i < dbLen; i++ {
		err = database.Change(db, "INSERT INTO "+c.DbTable+" VALUES("+fmt.Sprintf("%d", dbCandles.Array[i].Time)+
			", "+fmt.Sprintf("%f", dbCandles.Array[i].Open)+", "+fmt.Sprintf("%f", dbCandles.Array[i].Close)+", "+
			fmt.Sprintf("%f", dbCandles.Array[i].High)+", "+fmt.Sprintf("%f", dbCandles.Array[i].Low)+", "+
			fmt.Sprintf("%f", dbCandles.Array[i].Volume)+");")
		if err != nil {
			return err
		}
	}
	return nil
}
