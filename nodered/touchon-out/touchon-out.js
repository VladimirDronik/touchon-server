module.exports = function (RED) {
    class TouchOn_Out {
        constructor(config) {
            RED.nodes.createNode(this, config);
            this.targettype = config.targettype;
            this.targetid = Number(config.targetid);
            this.commandname = config.commandname;
            this.args = JSON.parse(config.args);
            var node = this;

            let touchonServer = RED.nodes.getNode(config.ws);
            if (!touchonServer) {
                node.status({fill: "red", shape: "ring", text: "no server"});

                this.on('close', function (removed, done) {
                    done();
                });

                return
            }
            let ws = touchonServer.ws

            ws.on('open', function() {
                node.status({fill: "green", shape: "dot", text: "connected"});
            });

            node.status({fill: "red", shape: "ring", text: "disconnected"});

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
