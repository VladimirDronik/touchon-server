module.exports = function (RED) {
    class TouchOn_In {
        constructor(config) {
            RED.nodes.createNode(this, config);
            let node = this;

            this.targettype = config.targettype;
            this.targetid = Number(config.targetid);
            this.eventname = config.eventname;

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

            node.onWSMessageWrapper = function(event){ node.onWSMessage(event); }
            node.ws.addEventListener('message', node.onWSMessageWrapper);

            node.on('close', node.onClose);
        }

        onWSOpen() {
            this.status({fill: "green", shape: "dot", text: "connected"});
        }

        onWSClose() {
            this.status({fill: "red", shape: "ring", text: "disconnected"});
        }

        onWSMessage(event) {
            let msg = JSON.parse(event.data)

            // Фильтруем сообщения по targetType, targetID, eventName

            if (this.targettype !== undefined && this.targettype !== '' && this.targettype !== msg.payload.target_type) {
                // console.log('>>>', this.targettype, msg.payload.target_type)
                return
            }

            if (this.targetid !== undefined && this.targetid > 0 && this.targetid !== msg.payload.target_id) {
                // console.log('>>>', this.targetid, msg.payload.target_id)
                return
            }

            if (this.eventname !== undefined && this.eventname !== '' && this.eventname !== msg.payload.name) {
                // console.log('>>>', this.eventname, msg.payload.event_name)
                return
            }

            msg.payload.props = msg.payload.payload
            delete msg.payload.payload

            this.send(msg)
        }

        onClose(removed, done) {
            this.ws.removeEventListener('open', this.onWSOpenWrapper)
            this.ws.removeEventListener('close', this.onWSCloseWrapper)
            this.ws.removeEventListener('message', this.onWSMessageWrapper)
            done();
        }
    }

    RED.nodes.registerType("touchon-in", TouchOn_In);
}
