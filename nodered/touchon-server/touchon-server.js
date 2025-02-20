const rws = require("reconnecting-websocket");

module.exports = function (RED) {
    class TouchOnServer {
        constructor(config) {
            RED.nodes.createNode(this, config);
            let node = this;

            node.ws = new rws(`ws://${config.host}:${config.port}${config.path}`);
            node.on('close', node.onClose);
        }

        addEventListener(eventName, handler) {
            this.ws.addEventListener(eventName, handler)
        }

        removeEventListener(eventName, handler) {
            this.ws.removeEventListener(eventName, handler)
        }

        onClose(removed, done) {
            this.ws.addEventListener('close', function(){ done(); });
            this.ws.close()
        }
    }

    RED.nodes.registerType("touchon-server", TouchOnServer);
}
