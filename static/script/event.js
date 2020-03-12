
function ClientButtonEvent(){
    astilectron.sendMessage("hello", function(message) {
        console.log("received " + message)
    });
}
