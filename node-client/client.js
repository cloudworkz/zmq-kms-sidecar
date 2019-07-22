"use strict";

const assert = require("assert");
const ZMQKMS = require("./zmqkms.js");
const zmqkms = new ZMQKMS();
zmqkms.on("error", console.error);

const call = async (i, log = false, decrypt = true) => {

    const cryptoKeyId = "super-secret-key"
    const plaintext = `3384395a-0dd4-491a-b0c1-f29e3f330933-${i}`;
    if (log) {
        console.log("plaintext:\n" + plaintext + "\n");
    }

    const cipher = await zmqkms.encrypt(cryptoKeyId, plaintext);
    if (log) {
        console.log("encrypted:\n" + cipher + "\n");
    }

    if (decrypt) {

        const decipher = await zmqkms.decrypt(cryptoKeyId, cipher);
        if (log) {
            console.log("decrypted:\n" + decipher + "\n");
        }

        assert.equal(plaintext, decipher);
    }
};

(async () => {

    // test
    await call("x", true);

    const startT = Date.now();
    const calls = [];
    for (let i = 0; i < 100000; i++) {
        // we are not making any decrypt calls to the google api
        // otherwise this will become an expensive benchmark
        calls.push(call(i, false, false));
    }

    await Promise.all(calls).then(() => {
        const tookMs = Date.now() - startT;
        console.log("Done.. took:", tookMs, "ms");
        zmqkms.close();
    });
})().catch(console.error);