package main

import (
	"fmt"
	"os"

	zmq "github.com/pebbe/zmq4"
)

func removeIndex(s []byte, index int) []byte {
	return append(s[:index], s[index+1:]...)
}

func main() {

	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	responder.Bind("tcp://*:5560")
	fmt.Println("Listening on port 5560.")

	args := os.Args
	if len(args) != 2 {
		panic("Requires exactly one argument: Config path.")
	}

	filePath := args[1]
	config := readConfig(filePath)
	/*
		{
			"projectID": "",
			"keyRingID": "",
			"locationID": "europe-west1",
			"cryptoKeyID": "",
			"cryptoKeyVersion": "1"
		}
	*/

	cryptoKeyName := getKeyName(config["projectID"], config["keyRingID"],
		config["locationID"], config["cryptoKeyID"], config["cryptoKeyVersion"])
	fmt.Println(cryptoKeyName)

	publicKey, err := getAsymmetricPublicKey(cryptoKeyName)
	if err != nil {
		panic(fmt.Sprintf("Failed to get public key: [%s]\n", err))
	}
	fmt.Println("Public key loaded.")

	berrorResponse := []byte("err")

	for {
		brequest, _ := responder.RecvBytes(0)
		hbyte := brequest[0]
		brequest = removeIndex(brequest, 0)

		if hbyte == byte(0) {
			cipher, err := encryptRSA(publicKey, brequest)
			if err != nil {
				fmt.Printf("Failed to encrypt: [%s]\n", err)
				responder.SendBytes(berrorResponse, 0)
			} else {
				responder.SendBytes(cipher, 0)
			}
		} else if hbyte == byte(1) {
			decCipher, err := decryptRSA(cryptoKeyName, brequest)
			if err != nil {
				fmt.Printf("Failed to encrypt: [%s]\n", err)
				responder.SendBytes(berrorResponse, 0)
			} else {
				responder.SendBytes(decCipher, 0)
			}
		} else {
			fmt.Printf("Unsupported operation: [%b]\n", hbyte)
			responder.SendBytes(berrorResponse, 0)
		}
	}
}
