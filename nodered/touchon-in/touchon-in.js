const rws = require("reconnecting-websocket");

module.exports = function (RED) {
    class TouchOn_In {
        constructor(config) {
            RED.nodes.createNode(this, config);
            this.targetType = config.targetType;
            this.targetID = config.targetID;
            this.eventName = config.eventName;
            var node = this;

            console.log(node);

            const ws = new rws('ws://localhost:8081/nodered');
            ws.on = ws.addEventListener

            // ws.addEventListener('error', console.error);

            ws.on('open', () => {
                // https://nodered.org/docs/creating-nodes/status
                // red, green, yellow, blue or grey
                node.status({fill: "green", shape: "dot", text: "connected"});
            });

            ws.on('message', (msg) => {
                // TODO фильтруем сообщения по targetType, targetID, eventName !!
                node.send(msg)
            });

            // ===========================================

            node.status({fill: "red", shape: "ring", text: "disconnected"});

            // node.send([[out1msg1, out1msg2], [out2msg1, out2msg2]]);
            // node.send(msg);
            // node.error('error msg');

            this.on('close', function (removed, done) {
                ws.on('close', function open() {
                    node.status({fill: "red", shape: "ring", text: "disconnected"});
                    done();
                });

                ws.close()
            });
        }
    }

    RED.nodes.registerType("touchon-in", TouchOn_In, {
        // settings: {
        //     touchonInTargetType: {
        //         value: "",
        //         exportable: true
        //     },
        //     touchonInTargetID: {
        //         value: 0,
        //         exportable: true
        //     },
        //     touchonInEventName: {
        //         value: "",
        //         exportable: true
        //     }
        // }
    });
}
