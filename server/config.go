package main

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	Host       string      `json:"host"`
	ProjectID  string      `json:"projectID"`
	KeyRingID  string      `json:"keyRingID"`
	LocationID string      `json:"locationID"`
	CryptoKeys []cryptoKey `json:"cryptoKeys"`
}

type cryptoKey struct {
	ID      string `json:"id"`
	Version string `json:"version"`
}

func readConfig(filePath string) config {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	cfg := config{}
	err = json.Unmarshal([]byte(file), &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
