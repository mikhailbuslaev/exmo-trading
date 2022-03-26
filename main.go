package main
import (
	"mikhailbuslaev/exmo/app/utils"
	"mikhailbuslaev/exmo/app/query"
	"encoding/json"
	"fmt"
)
func main() {
	user := utils.LoadUser()
	fmt.Println(user.PublicKey)
	resp, err := query.Do("user_info", nil, user)
	utils.PrintResponse(resp, err)

	jsonresp, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Log encoding fail")
	}
	utils.RecordLog(jsonresp)
}