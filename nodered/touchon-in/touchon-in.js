module.exports = function (RED) {
    class TouchOn_In {
        constructor(config) {
            RED.nodes.createNode(this, config);
            this.targettype = config.targettype;
            this.targetid = config.targetid;
            this.eventname = config.eventname;
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

            ws.on('message', function(event) {
                let msg = JSON.parse(event.data)

                // Фильтруем сообщения по targetType, targetID, eventName

                if (node.targettype !== undefined && node.targettype !== '' && node.targettype !== msg.payload.target_type) {
                    // console.log('>>>', node.targettype, msg.payload.target_type)
                    return
                }

                const targetID = Number(node.targetid)
                if (node.targetid !== undefined && targetID > 0 && targetID !== msg.payload.target_id) {
                    // console.log('>>>', node.targetid, msg.payload.target_id)
                    return
                }

                if (node.eventname !== undefined && node.eventname !== '' && node.eventname !== msg.payload.name) {
                    // console.log('>>>', node.eventname, msg.payload.event_name)
                    return
                }

                node.send(msg)
            });

            node.status({fill: "red", shape: "ring", text: "disconnected"});

            this.on('close', function (removed, done) {
                ws.on('close', function () {
                    node.status({fill: "red", shape: "ring", text: "disconnected"});
                    done();
                });

                ws.close()
            });
        }
    }

    RED.nodes.registerType("touchon-in", TouchOn_In);
}
