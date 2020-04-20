const Events = {
    appHostStart            : "app.event.host",
    appWorkerStart          : "app.worker.start",
    appGenerateToken        : "app.generate.token",
    windowDeviceStatus      : "window.device.status",
    windowNetworkStatus     : "window.network.status",
    windowTaskProgress      : "window.task.progress",
    windowSendToken         : "window.send.token"
};

let websocket = undefined;

const connect = () => {
    if (websocket === undefined)
        websocket = new WebSocket("ws://localhost:" + global.backendPort + "/web/app/events");
}

const sendMessage = (event, message) => {
    if (websocket === undefined) {
        console.log("Error Send Message not connected app");
        return;
    };
    websocket.send(JSON.stringify({
        "event": event,
        "message": message,
    }));
}

const receiveMessage = (event, callback) => {
    if (websocket === undefined) {
        console.log("Error Send Message not connected app");
        return;
    };

    websocket.addEventListener('message', function (message) {
        const json = JSON.parse(message.data);
        if (json.event.e === event) {
            callback(obj.message);
        }
    });
}

export {
    Events,
    connect,
    sendMessage,
    receiveMessage
};