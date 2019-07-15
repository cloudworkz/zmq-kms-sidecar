"use strict";

const ZMQ = require("zeromq");

const ENCRYPT_BUFFER = Buffer.from([0]);
const DECRYPT_BUFFER = Buffer.from([1]);

const encrypt = (socket, plaintextStr, encoding = "utf8") => {
    return new Promise((resolve, reject) => {
        const pt = Buffer.concat([ENCRYPT_BUFFER, Buffer.from(plaintextStr, encoding)]);
        socket.once("message", (message) => {

            if (message.length <= 3) {
                return reject(new Error("Failed to encrypt."));
            }

            return resolve(message.toString("hex"));
        });
        socket.send(pt, 0);
    });
};

const decrypt = (socket, cipherStr, encoding = "hex") => {
    return new Promise((resolve, reject) => {
        const ct = Buffer.concat([DECRYPT_BUFFER, Buffer.from(cipherStr, encoding)]);
        socket.once("message", (message) => {

            if (message.length <= 3) {
                return reject(new Error("Failed to encrypt."));
            }

            return resolve(message.toString("utf8"));
        });
        socket.send(ct, 0);
    });
};

const socket = ZMQ.socket("req");
socket.connect("tcp://127.0.0.1:5560");

const call1 = (async () => {
    const plaintext = "3384395a-0dd4-491a-b0c1-f29e3f330933-a";
    console.log("plaintext:\n" + plaintext + "\n");

    const ctT = Date.now();
    const cipher = await encrypt(socket, plaintext);
    const cteT = Date.now();

    const dtT = Date.now();
    const decipher = await decrypt(socket, cipher);
    const dteT = Date.now();

    console.log("encrypted:\n" + cipher + "\n", (cteT - ctT));
    console.log("decrypted:\n" + decipher + "\n", (dteT - dtT));
    console.log("eq:", plaintext === decipher);
})();

const call2 = (async () => {
    const plaintext = "3384395a-0dd4-491a-b0c1-f29e3f330933-b";
    console.log("plaintext:\n" + plaintext + "\n");

    const ctT = Date.now();
    const cipher = await encrypt(socket, plaintext);
    const cteT = Date.now();

    const dtT = Date.now();
    const decipher = await decrypt(socket, cipher);
    const dteT = Date.now();

    console.log("encrypted:\n" + cipher + "\n", (cteT - ctT));
    console.log("decrypted:\n" + decipher + "\n", (dteT - dtT));
    console.log("eq:", plaintext === decipher);
})();

const call3 = (async () => {
    const plaintext = "3384395a-0dd4-491a-b0c1-f29e3f330933-c";
    console.log("plaintext:\n" + plaintext + "\n");

    const ctT = Date.now();
    const cipher = await encrypt(socket, plaintext);
    const cteT = Date.now();

    const dtT = Date.now();
    const decipher = await decrypt(socket, cipher);
    const dteT = Date.now();

    console.log("encrypted:\n" + cipher + "\n", (cteT - ctT));
    console.log("decrypted:\n" + decipher + "\n", (dteT - dtT));
    console.log("eq:", plaintext === decipher);
})();

Promise.all([
    call1,
    call2,
    call3
]).then(() => {
    socket.close();
}).catch(console.log);