package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

func Connect(dbConfigName string) (*sql.DB, error) {
	file, err := os.Open(dbConfigName)

	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("Successfully opened database-config.json")
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	var config Config
	json.Unmarshal(byteValue, &config)

	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", config.Host,
		config.Port, config.User, config.Password, config.Dbname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Database opening fail")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Database connection fail")
	} else {
		fmt.Println("Database connection is successful")
	}
	return db, err
}

func Change(db *sql.DB, request string) error {
	_, err := db.Exec(request)
	if err != nil {
		return err
	}
	return nil
}

func Select(db *sql.DB, request string) (*sql.Rows, error) {
	result, err := db.Query(request)

	if err != nil {
		fmt.Println("Select query failure")
	} else {
		fmt.Println("Select query correct")
	}
	return result, err
}
