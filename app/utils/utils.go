package utils

import (
	"fmt"
	"io/ioutil"
	. "mikhailbuslaev/exmo/app/types"
	"os"

	"gopkg.in/yaml.v2"
)

func PrintResponse(resp map[string]interface{}, err error) {
	if err != nil {
		fmt.Println("Error while do request")
	} else {
		fmt.Println(resp)
	}
}

func Record(data []byte, way string) {

	file, err := os.OpenFile(way, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error while writing logs")
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error while writing logs")
	}
}

func RecordNewLine(way string) {
	file, err := os.OpenFile(way, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error while writing logs")
	}
	defer file.Close()
	_, err = file.Write([]byte("\n"))
	if err != nil {
		fmt.Println("Error while writing logs")
	}
}


func Rewrite(data []byte, way string) {
	err := os.Truncate(way, 0)
	if err != nil {
		fmt.Println("Error while cleaning file")
	}
	file, err := os.OpenFile(way, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		fmt.Println("Error while open file")
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error while rewriting file")
	}
}

func LoadUser() *User {
	user := &User{}
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
