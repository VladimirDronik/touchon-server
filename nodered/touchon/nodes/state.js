module.exports = function (RED) {
    class TouchOn_State {
        constructor(config) {
            RED.nodes.createNode(this, config);
            let node = this;

            node.targettype = config.targettype;
            node.targetid = Number(config.targetid);

            node.ws = RED.nodes.getNode(config.ws);
            if (!node.ws) {
                node.status({fill: "red", shape: "ring", text: "no server"});
                node.on('close', function (removed, done) { done(); });
                return
            }

            node.status({fill: "red", shape: "ring", text: "disconnected"});

            node.onWSOpenWrapper = function(){ node.onWSOpen(); }
            node.ws.addEventListener('open', node.onWSOpenWrapper);

            node.onWSCloseWrapper = function(){ node.onWSClose(); }
            node.ws.addEventListener('close', node.onWSCloseWrapper);

            node.on('input', node.onInput);
            node.on('close', node.onClose);
        }

        onWSOpen() {
            this.status({fill: "green", shape: "dot", text: "connected"});
        }

        onWSClose() {
            this.status({fill: "red", shape: "ring", text: "disconnected"});
        }

        onInput(oldMsg, send, done) {
            var node = this;
            let uri = `http://${node.ws.host}:${node.ws.port}/${node.targettype}s/${node.targetid}/state`

            fetch(uri)
                .then((response) => {
                    if (!response.ok) {
                        throw new Error(response.statusText);
                    }

                    let msg = response.json();
                    if (msg.error) {
                        throw new Error(msg.error);
                    }

                    return msg
                })
                .then((msg) => {
                    let m = {
                        oldMessage: oldMsg, // Прокидываем исходное сообщение дальше
                        payload: msg.response
                    }
                    m.payload.props = m.payload.payload
                    delete m.payload.payload
                    send(m)
                })
                .catch((error) => node.status({fill: "red", shape: "ring", text: error}))
                .finally(() => done())
        }

        onClose(removed, done) {
            this.ws.removeEventListener('open', this.onWSOpenWrapper)
            this.ws.removeEventListener('close', this.onWSCloseWrapper)
            done();
        }
    }

    RED.nodes.registerType("touchon-state", TouchOn_State);
}
