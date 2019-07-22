# zmq-kms-sidecar

# What?

This is a sidecar that helps you to work with Google Cloud KMS from Node.js

# How?

You basically run a small go server that communicates with your node app using
a small byte based protocol on top of the zeromq request and response pattern.

# Why?

As of now, Node.js does not correctly support the required OEP padding.

# Trying it out locally

## The server (sidecar)

You will have to add a `config.json` file to `./server/config.json`
and fill out the following content:

```javascript
{
    "host": "tcp://*:5560",
    "projectID": "your-gcp-project",
    "keyRingID": "your-key-ring-id",
    "locationID": "europe-west1",
    "cryptoKeys": [
        {
            "cryptoKeyID": "your-key-name-1",
            "cryptoKeyVersion": "1",
        },
        {
            "cryptoKeyID": "your-key-name-2",
            "cryptoKeyVersion": "4",
        }
    ]
}
```

You can then run the following to compile and start the server:
(Please Note: This requires `Go`, `pkg-config` and `zmq` being installed on your computer)

```bash
cd ./server
go get .
go build .
./zmq-kms ./config.json
```

## The client

Head into `cd ./node-client` and run:
(Please Note: This requires`Node.js` and `yarn` being installed on your computer)

```bash
yarn
yarn start
```

# Using this in production

Its quite easy to wrap the server component in a Docker container.
All you need is the compiled binary and your JSON config file.

Implementing the client is also quite easy, you will just have to add the dependencies `zeromq`
and `uuid` to your project. And copy the files `./node-client/zmqdr.js` (a small wrapper around zmq that adds
callbacks to the messages sends via call-stack) and `./node-client/zmqkms.js` (an even smaller wrapper around zmqdr
that gives you a simple encrypt/decrypt promisfied interface). You can pass the connection string to the constructor
of zmqdr or zmqkms. `./node-client/client.js` gives you a starting point.
