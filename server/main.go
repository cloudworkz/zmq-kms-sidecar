package main

import (
	"encoding/hex"
	"encoding/json"
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

	responder, _ := zmq.NewSocket(zmq.REP)
	defer responder.Close()
	responder.Bind(config.Host)
	fmt.Println("Listening on at: " + config.Host)

	publicKeys := make(map[string]interface{}, len(config.CryptoKeys))
	for _, cryptoKey := range config.CryptoKeys {
		cryptoKeyName := getKeyName(config.ProjectID, config.KeyRingID,
			config.LocationID, cryptoKey.ID, cryptoKey.Version)
		fmt.Println(cryptoKeyName)

		publicKey, err := getAsymmetricPublicKey(cryptoKeyName)
		if err != nil {
			panic(fmt.Sprintf("Failed to get public key: [%s]\n", err))
		}
		fmt.Printf("Public key '%v' (version '%v') loaded.\n", cryptoKey.ID, cryptoKey.Version)
		publicKeys[cryptoKey.ID] = publicKey
	}

	berrorResponse := []byte("err")

	for {
		brequest, _ := responder.RecvBytes(0)
		hbyte := brequest[0]     // first message byte is the command header
		idbyte := brequest[1:37] // next 36 bytes are the uuidv4 message identifier
		brequest = brequest[37:] // rest of the message is a stringified JSON payload

		// fmt.Printf("Header: [%b]\n", hbyte)
		// fmt.Printf("ID: [%s]\n", string(idbyte))
		// fmt.Printf("Payload: [%s]\n", hex.EncodeToString(brequest))

		if hbyte == byte(0) {
			request := encryptRequest{}
			err := json.Unmarshal(brequest, &request)
			if err != nil {
				fmt.Printf("Failed to unmarshal encrypt request: [%b]\n", hbyte)
				responder.SendBytes(append(idbyte, berrorResponse...), 0)
				continue
			}

			publicKey, ok := publicKeys[request.CryptoKeyID]
			if !ok {
				fmt.Printf("Cryptokey with id '%v' not found\n", request.CryptoKeyID)
				responder.SendBytes(append(idbyte, berrorResponse...), 0)
				continue
			}

			cipher, err := encryptRSA(publicKey, []byte(request.Plaintext))
			if err != nil {
				fmt.Printf("Failed to encrypt: [%s]\n", err)
				responder.SendBytes(append(idbyte, berrorResponse...), 0)
				continue
			}

			responder.SendBytes(append(idbyte, cipher...), 0)
		} else if hbyte == byte(1) {
			request := decryptRequest{}
			err := json.Unmarshal(brequest, &request)
			if err != nil {
				fmt.Printf("Failed to unmarshal decrypt request: [%b]\n", hbyte)
				responder.SendBytes(append(idbyte, berrorResponse...), 0)
				continue
			}

			var version string
			for _, key := range config.CryptoKeys {
				if key.ID == request.CryptoKeyID {
					version = key.Version
					break
				}
			}
			if version == "" {
				fmt.Printf("Cryptokey with id '%v' not found in config\n", request.CryptoKeyID)
				responder.SendBytes(append(idbyte, berrorResponse...), 0)
				continue
			}

			cipherBytes, err := hex.DecodeString(request.Cipher)
			if err != nil {
				fmt.Printf("Failed to decode hex string: [%b]\n", hbyte)
				responder.SendBytes(append(idbyte, berrorResponse...), 0)
				continue
			}

			cryptoKeyName := getKeyName(config.ProjectID, config.KeyRingID,
				config.LocationID, request.CryptoKeyID, version)
			decCipher, err := decryptRSA(cryptoKeyName, cipherBytes)
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

// encryptRequest is the JSON paylaod we receive via ZeroMQ
type encryptRequest struct {
	CryptoKeyID string `json:"cryptoKeyID"`
	Plaintext   string `json:"plaintext"`
}

type decryptRequest struct {
	CryptoKeyID string `json:"cryptoKeyID"`
	Cipher      string `json:"cipher"`
}
