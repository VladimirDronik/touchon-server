const rws = require("reconnecting-websocket");

module.exports = function (RED) {
    class TouchOnServer {
        constructor(config) {
            RED.nodes.createNode(this, config);
            this.ws = new rws(`ws://${config.host}:${config.port}${config.path}`);
            this.ws.on = this.ws.addEventListener

            this.on('close', function (removed, done) {
                done();
            });
        }
    }

    RED.nodes.registerType("touchon-server", TouchOnServer);
}
