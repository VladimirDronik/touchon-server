module.exports = function (RED) {
    class TouchOn_Out {
        constructor(config) {
            RED.nodes.createNode(this, config);
            let node = this;

            node.targettype = config.targettype;
            node.targetid = Number(config.targetid);
            node.commandname = config.commandname;
            node.args = JSON.parse(config.args);

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

        onInput(mess, send, done) {
            let msg = {
                type: 'command',
                target_type: this.targettype,
                target_id: this.targetid,
                name: this.commandname,
                payload: this.args,
            }

            this.ws.send(JSON.stringify(msg))
            done()
        }

        onClose(removed, done) {
            this.ws.removeEventListener('open', this.onWSOpenWrapper)
            this.ws.removeEventListener('close', this.onWSCloseWrapper)
            done();
        }
    }

    RED.nodes.registerType("touchon-out", TouchOn_Out);
}
