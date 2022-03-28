package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mikhailbuslaev/exmo/app/query"
	"mikhailbuslaev/exmo/app/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func ListenHandler(w http.ResponseWriter, r *http.Request) {
	user := utils.LoadUser()
	resp, err := query.Do("user_info", nil, user)
	if err != nil {
		fmt.Println("Log encoding fail")
	}
	//	utils.PrintResponse(resp, err)

	jsonresp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Log encoding fail")
	}
	utils.RecordNewLine("logs/test.log")
	utils.Record(jsonresp, "logs/test.log")

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/listen", ListenHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:1111",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
