"use strict";

const EventEmitter = require("events");
const ZMQDR = require("./zmqdr.js");

const ENCRYPT_BUFFER = Buffer.from([0]);
const DECRYPT_BUFFER = Buffer.from([1]);

class ZMQKMS extends EventEmitter {

    /**
     * Creates new instance of wrapper around zmqdr
     * @param {string} conStr - e.g. "tcp://127.0.0.1:5560"
     * @param {number} maxStackSize - default is 0, if set there should not be more parallel calls
     */
    constructor(conStr = "tcp://127.0.0.1:5560", maxStackSize = 0) {
        super();
        this.socket = new ZMQDR(conStr);
        this.socket.on("error", (error) => this.emit("error", error));
    }

    encrypt(plaintextStr, encoding = "utf8") {
        return new Promise((resolve, reject) => {
            this.socket.send(ENCRYPT_BUFFER, Buffer.from(plaintextStr, encoding), (error, message) => {

                if (error) {
                    return reject(error);
                }

                if (message.length <= 3) {
                    return reject(new Error("Failed to encrypt."));
                }

                return resolve(message.toString("hex"));
            });
        });
    }

    decrypt(cipherStr, encoding = "hex") {
        return new Promise((resolve, reject) => {
            this.socket.send(DECRYPT_BUFFER, Buffer.from(cipherStr, encoding), (error, message) => {

                if (error) {
                    return reject(error);
                }

                if (message.length <= 3) {
                    return reject(new Error("Failed to decrypt."));
                }

                return resolve(message.toString("utf8"));
            });
        });
    }

    close() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }
}

module.exports = ZMQKMS;
