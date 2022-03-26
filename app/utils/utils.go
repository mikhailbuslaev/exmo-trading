package utils

import (
	"fmt"
	"os"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	. "mikhailbuslaev/exmo/app/types"
)

func PrintResponse(resp map[string]interface{}, err error) {
	if err != nil {
		fmt.Println(resp)
	} else {
		fmt.Println("Error while do request")
	}
}

func RecordLog(data []byte) {

	file, err := os.OpenFile("logs/test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Error while writing logs")
    }
	
	
   _, err3 := file.Write(data) 
    if err3 != nil {
		fmt.Println("Error while writing logs")
    }

	file.Close()
}

func LoadUser() *User{
	var user *User
	file, err := ioutil.ReadFile("configs/user-config.yaml")
    if err != nil {
        fmt.Println("Error while open user-config.yaml")
    }

	err = yaml.Unmarshal(file, user)
    if err != nil {
        fmt.Println("Error while unmarshalling user-config.yaml")
    }

    return user
}
