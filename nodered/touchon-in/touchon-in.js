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

            // ws.addEventListener('open', () => {
            //     ws.send('hello!');
            // });

            // ws.addEventListener('error', console.error);

            // ws.addEventListener('open', function open() {
            //     ws.send('something');
            // });

            ws.on('message', function message(data) {
                console.log('received: %s', data);
                ws.send(data)
            });

            // ===========================================

            node.status({fill: "red", shape: "ring", text: "disconnected"});

            node.on('input', function (msg, send, done) {
                msg.payload = msg.payload.toLowerCase();
                send(msg);

                // https://nodered.org/docs/creating-nodes/status
                // red, green, yellow, blue or grey
                node.status({fill: "green", shape: "dot", text: "connected"});
                done(); // done('error msg')
            });

            // node.send([[out1msg1, out1msg2], [out2msg1, out2msg2]]);
            // node.send(msg);
            // node.error('error msg');

            this.on('close', function (removed, done) {
                // doSomethingWithACallback(function() { done(); });
                node.status({fill: "red", shape: "ring", text: "disconnected"});
                done();
            });
        }
    }

    RED.nodes.registerType("touchon-in", TouchOn_In, {
        settings: {
            touchonInTargetType: {
                value: "",
                exportable: true
            },
            touchonInTargetID: {
                value: 0,
                exportable: true
            },
            touchonInEventName: {
                value: "",
                exportable: true
            }
        }
    });
}
