"use strict";

const EventEmitter = require("events");
const ZMQ = require("zeromq");
const uuid = require("uuid");

const EMPTY_BUFFER = Buffer.from([]);

class ZMQDR extends EventEmitter {

    /**
     * Creates new instance of wrapper around zmq socket
     * @param {string} conStr - e.g. "tcp://127.0.0.1:5560"
     * @param {number} maxStackSize - default is 0, if set there should not be more parallel calls
     */
    constructor(conStr, maxStackSize = 0) {
        super();

        this.maxStackSize = maxStackSize;
        this.socket = ZMQ.socket("req");
        this.socket.connect(conStr);
        this.socket.on("message", this._onMessage.bind(this));
        this._stack = {};
    }

    _onMessage(message) {

        let identifier = null;
        try {
            identifier = message.slice(0, 36).toString("utf8");
        } catch (error) {
            // empty
        }

        if (!identifier) {
            return this.emit("error", new Error("Failed to parse identifier from message: " + message.toString("hex")));
        }

        if (!this._stack[identifier]) {
            return this.emit("error", new Error("No identifier present in stack for: " + identifier));
        }

        const payload = message.slice(36, message.length);

        this._stack[identifier](null, payload);
        delete this._stack[identifier];
    }

    send(headBuffer, messageBuffer, callback) {

        if (this.maxStackSize > 0 && this.maxStackSize >= this.getStackSize()) {
            this.emit("error", new Error("Max stack size exceeded: " + this.maxStackSize));
            Object.keys(this._stack).forEach((identifier) => {
                this._stack[identifier](new Error("Stack size exceeded, removed this call."));
                delete this._stack[identifier];
            });
        }

        if (!headBuffer) {
            headBuffer = EMPTY_BUFFER;
        }

        const id = uuid.v4();
        const pt = Buffer.concat([headBuffer, Buffer.from(id, "utf8"), messageBuffer]);
        this._stack[id] = callback;
        this.socket.send(pt, 0);
        return id;
    }

    getStackSize() {
        return Object.keys(this._stack).length;
    }

    close() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }
}

module.exports = ZMQDR;
