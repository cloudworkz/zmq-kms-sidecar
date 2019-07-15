package main

import (
	"fmt"

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

	projectID := "rd-bigdata-int-v002"
	keyRingID := "eden"
	locationID := "europe-west1"
	cryptoKeyID := "eden-reference"
	cryptoKeyVersion := "1"
	cryptoKeyName := getKeyName(projectID, keyRingID, locationID, cryptoKeyID, cryptoKeyVersion)

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
