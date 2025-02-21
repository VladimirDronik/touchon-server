module.exports = function (RED) {
    "use strict";

    function CopyNode(config) {
        RED.nodes.createNode(this, config);
        var node = this;

        node.outputs = config.outputs;

        this.on('input', function (msg, send, done) {
            var ports = new Array(node.outputs).fill(null);
            for (var i = 0; i < node.outputs; i++) {
                var newMsg = RED.util.cloneMessage(msg);
                newMsg._msgid = RED.util.generateId();
                ports[i] = newMsg;
            }

            send(ports);
            done();
        });
    }

    RED.nodes.registerType("copy", CopyNode);
}