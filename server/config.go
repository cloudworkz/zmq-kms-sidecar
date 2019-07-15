package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func readConfig(filePath string) map[string]string {

	jsonFile, err := os.Open(filePath)
	var result map[string]string
	if err != nil {
		fmt.Println(err)
		return result
	}
	fmt.Println("Successfully opened " + filePath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &result)
	return result
}
