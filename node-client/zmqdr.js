"use strict";

const ZMQ = require("zeromq");
const uuid = require("uuid").v4;

class ZMQDR {

    constructor(conStr = "tcp://127.0.0.1:5560") {
        this.socket = ZMQ.socket("req");
        this.socket.connect(conStr);
        this.socket.on("message", this._onMessage.bind(this));
        this.stack = {};
    }

    _onMessage(message) {

    }

}