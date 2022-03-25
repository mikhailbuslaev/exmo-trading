package utils

import (

	"fmt"
	"os"
)

func PrintResponse(resp map[string]interface{}, err error) {
	if err != nil {
		fmt.Println(resp)
	} else {
		fmt.Printf("Error: %s\n", err.Error())
	}
}

func RecordLog(data []byte) {

	file, err := os.Open("logs/requests.txt")
    if err != nil {
        fmt.Printf("Error: %s\n", err.Error())
    }
    defer file.Close()
    _, err3 := file.WriteAt(data, 0)
    if err3 != nil {
		fmt.Printf("Error: %s\n", err.Error())
    }
}
