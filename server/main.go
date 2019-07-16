package main

import (
	"fmt"
	"os"

	zmq "github.com/pebbe/zmq4"
)

func main() {

	args := os.Args
	if len(args) != 2 {
		panic("Requires exactly one argument: Config path.")
	}

	filePath := args[1]
	config := readConfig(filePath)
	/*
		{
			"host": "tcp://*:5560",
			"projectID": "",
			"keyRingID": "",
			"locationID": "europe-west1",
			"cryptoKeyID": "",
			"cryptoKeyVersion": "1"
		}
	*/

	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	responder.Bind(config["host"])
	fmt.Println("Listening on at: " + config["host"])

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
		hbyte := brequest[0]     // first message byte is the command header
		idbyte := brequest[1:37] // next 36 bytes are the uuidv4 message identifier
		brequest = brequest[37:] // rest of the message is the payload

		// fmt.Printf("Header: [%b]\n", hbyte)
		// fmt.Printf("ID: [%s]\n", string(idbyte))
		// fmt.Printf("Payload: [%s]\n", hex.EncodeToString(brequest))

		if hbyte == byte(0) {
			cipher, err := encryptRSA(publicKey, brequest)
			if err != nil {
				fmt.Printf("Failed to encrypt: [%s]\n", err)
				responder.SendBytes(append(idbyte, berrorResponse...), 0)
			} else {
				responder.SendBytes(append(idbyte, cipher...), 0)
			}
		} else if hbyte == byte(1) {
			decCipher, err := decryptRSA(cryptoKeyName, brequest)
			if err != nil {
				fmt.Printf("Failed to decrypt: [%s]\n", err)
				responder.SendBytes(append(idbyte, berrorResponse...), 0)
			} else {
				responder.SendBytes(append(idbyte, decCipher...), 0)
			}
		} else {
			fmt.Printf("Unsupported operation: [%b]\n", hbyte)
			responder.SendBytes(append(idbyte, berrorResponse...), 0)
		}
	}
}
