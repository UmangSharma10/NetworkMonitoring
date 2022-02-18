package main

import (
	"encoding/base64"
	json "encoding/json"
	"fmt"
	"os"
	"sample/plugins"
)

func main() {

	channel := make(chan string)
	recevedARG1 := os.Args[1]

	jsonDecodedString, err := base64.StdEncoding.DecodeString(recevedARG1)

	if err != nil {
		panic(err)
	}

	var credMap map[string]string
	//fmt.Println(string(jsonDecodedString))
	json.Unmarshal(jsonDecodedString, &credMap)
	//fmt.Println(credMap)
	fmt.Println("it will collect addresss")

	if string(credMap["device"]) == "windows" {
		go plugins.CollectAddress(credMap, channel)

	} else if string(credMap["device"]) == "linux" {
		go plugins.CollectSSH(credMap, channel)
	}

	for i := 0; i < len(credMap); i++ {
		result := <-channel
		fmt.Println(result)
	}
	close(channel)

	/*for _, v := range credMap {
	if string(v["device"]) == "windows" {
		go plugins.CollectAddress(v, channel)

	} else if string(v["device"]) == "linux" {
		go plugins.CollectSSH(v, channel)
	}*/

	//}
	//println(len(credMap))

}
