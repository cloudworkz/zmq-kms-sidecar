package main

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

func main() {

	requester, _ := zmq.NewSocket(zmq.REQ)
	defer requester.Close()
	requester.Connect("tcp://localhost:5560")

	bencryptCmd := []byte{0}
	bdecryptCmd := []byte{1}

	bplaintext := append(bencryptCmd, []byte("3384395a-0dd4-491a-b0c1-f29e3f330933")...)

	requester.SendBytes(bplaintext, 0)
	cipherReply, _ := requester.RecvBytes(0)
	fmt.Printf("Received (encryption) reply %s\n", string(cipherReply))

	bcipher := append(bdecryptCmd, cipherReply...)

	requester.SendBytes(bcipher, 0)
	decipherReply, _ := requester.RecvBytes(0)
	fmt.Printf("Received (decryption) reply %s\n", string(decipherReply))
}
