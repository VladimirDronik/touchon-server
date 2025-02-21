const rws = require("reconnecting-websocket");

module.exports = function (RED) {
    class TouchOnServer {
        constructor(config) {
            RED.nodes.createNode(this, config);
            let node = this;
            node.host = config.host
            node.port = config.port
            node.path = config.path

            node.ws = new rws(`ws://${config.host}:${config.port}${config.path}`);
            node.on('close', node.onClose);
        }

        addEventListener(eventName, handler) {
            this.ws.addEventListener(eventName, handler)
        }

        removeEventListener(eventName, handler) {
            this.ws.removeEventListener(eventName, handler)
        }

        wsSend(data) {
            this.ws.send(data)
        }

        onClose(removed, done) {
            if (this.ws.readyState === this.ws.OPEN) {
                this.ws.addEventListener('close', function(){ done(); });
                this.ws.close()
            } else {
                done();
            }
        }
    }

    RED.nodes.registerType("touchon-server", TouchOnServer);
}
