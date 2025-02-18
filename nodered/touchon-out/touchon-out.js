const rws = require("reconnecting-websocket");

module.exports = function (RED) {
    class TouchOn_Out {
        constructor(config) {
            RED.nodes.createNode(this, config);
            this.targettype = config.targettype;
            this.targetid = Number(config.targetid);
            this.commandname = config.commandname;
            this.args = JSON.parse(config.args);
            var node = this;

            const ws = new rws('ws://localhost:8081/nodered');
            ws.on = ws.addEventListener

            // ws.addEventListener('error', console.error);

            ws.on('open', function() {
                // https://nodered.org/docs/creating-nodes/status
                // red, green, yellow, blue or grey
                node.status({fill: "green", shape: "dot", text: "connected"});
            });

            // ===========================================

            node.status({fill: "red", shape: "ring", text: "disconnected"});

            // node.send([[out1msg1, out1msg2], [out2msg1, out2msg2]]);
            // node.send(msg);
            // node.error('error msg');

            this.on('input', function (mess, send, done) {
                let msg = {
                    type: 'command',
                    target_type: node.targettype,
                    target_id: node.targetid,
                    name: node.commandname,
                    payload: node.args,
                }

                ws.send(JSON.stringify(msg))
                done()
            });

            this.on('close', function (removed, done) {
                ws.on('close', function () {
                    node.status({fill: "red", shape: "ring", text: "disconnected"});
                    done();
                });

                ws.close()
            });
        }
    }

    RED.nodes.registerType("touchon-out", TouchOn_Out);
}
