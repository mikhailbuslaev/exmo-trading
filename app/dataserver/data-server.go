package dataserver

import (
	"errors"
	"exmo-trading/app/data"
	"exmo-trading/app/database"
	"exmo-trading/app/query"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type Handler struct {
	Context Context
}

type Context struct {
	Symbol         string
	Resolution     string
	CandlesFile    string
	LastCandleFile string
	DbTable        string
}

func (h *Handler) LoadCandles(from, to string) error {
	q := query.GetQuery{Method: "candles_history?symbol=" + h.Context.Symbol +
		"&resolution=" + h.Context.Resolution + "&from=" + from + "&to=" + to}

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
	candles.Array = make([]data.Candle, 0, 1000)
	fmt.Println(len([]byte(body)))
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

		err = data.Rewrite(&candles.Array[len(candles.Array)-1], h.Context.LastCandleFile)
		if err != nil {
			fmt.Println("error while rewriting data")
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

func (h *Handler) UpdateCandles() error {
	candle := &data.Candle{}
	err := candle.Read(h.Context.LastCandleFile)
	if err != nil {
		return err
	}
	t := time.Now()
	err = h.LoadCandles(fmt.Sprintf("%d", (candle.Time/1000)+1), fmt.Sprintf("%d", t.Unix()))
	if err != nil {
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
	intresol, err := strconv.ParseInt(h.Context.Resolution, 10, 64)
	if err != nil {
		return err
	}
	err = h.LoadCandles(fmt.Sprintf("%d", t.Unix()-60*intresol*1000), fmt.Sprintf("%d", t.Unix()))
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) SyncDbAndCache(dbConfigName string) error {
	db, err := database.Connect(dbConfigName)
	if err != nil {
		return err
	}

	rows, err := database.Select(db, "SELECT * FROM "+h.Context.DbTable+" ORDER BY time ASC;")
	if err != nil {
		return err
	}

	cacheCandles := &data.Candles{}
	err = cacheCandles.Read(h.Context.CandlesFile)
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

	err = database.Change(db, "DELETE * FROM "+h.Context.DbTable+";")
	if err != nil {
		return err
	}

	for i := 0; i < dbLen; i++ {
		err = database.Change(db, "INSERT INTO "+h.Context.DbTable+" VALUES("+fmt.Sprintf("%d", dbCandles.Array[i].Time)+
			", "+fmt.Sprintf("%f", dbCandles.Array[i].Open)+", "+fmt.Sprintf("%f", dbCandles.Array[i].Close)+", "+
			fmt.Sprintf("%f", dbCandles.Array[i].High)+", "+fmt.Sprintf("%f", dbCandles.Array[i].Low)+", "+
			fmt.Sprintf("%f", dbCandles.Array[i].Volume)+");")
		if err != nil {
			return err
		}
	}
	return nil
}

func Launch() {

	fivemin := &Handler{
		Context: Context{
			Symbol:         "BTC_USD",
			Resolution:     "5",
			CandlesFile:    "cache/5min-candles.csv",
			LastCandleFile: "cache/5min-last-candle.csv",
			DbTable:        "5min-candles",
		},
	}

	err := fivemin.InitCandles()
	if err != nil {
		fmt.Println(err)
	}
	for {
		err := fivemin.UpdateCandles()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("waiting new candles")
		time.Sleep(60 * time.Second)
	}
}
