let websocket = undefined;

const connect = () => {
    if (websocket === undefined)
        websocket = new WebSocket("ws://localhost:" + global.backendPort + "/web/app/events");
}

const sendMessage = (event, message) => {
    console.log(websocket);
    if (websocket === undefined) {
        console.log("Error Send Message not connected app");
        return;
    };
    websocket.send(JSON.stringify({
        "event": event,
        "message": message,
    }));
}

const onMessage = () => {
    socket.addEventListener('message', function (event) {
        let obj = JSON.parse(message.data);

        // event name
        console.log(obj.event);

        // event data
        console.log(obj.AtrNameInFrontend);
    });
}

export {
    connect,
    sendMessage,
    onMessage
};