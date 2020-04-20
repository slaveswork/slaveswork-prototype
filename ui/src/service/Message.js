const Events = {
    appHostStart            : "app.event.host",
    appWorkerStart          : "app.worker.start",
    appGenerateToken        : "app.generate.token",
    appConnectDevice        : "app.connect.device",
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
    const sendData = JSON.stringify({
        "event": event,
        "message": message,
    });
    console.log(sendData);
    websocket.send(sendData);
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