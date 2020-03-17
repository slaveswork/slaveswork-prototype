
window.onload = function(){
    var start = document.getElementById("electron_start");
    console.log(start), start.addEventListener("astilectron-ready", (function () {
        astilectron.sendMessage("hello", (function (e) {
            console.log("received " + e)
        }))
    }))
};

document.getElementById("host").addEventListener("click", function() {
    // Send host 'start' message to GO
    astilectron.sendMessage("host", function(message) {
        // Token element have only class name. So I used that.
        var tokenElement = document.getElementsByClassName('token'); 
        tokenElement.item(0).innerText(message);
    });
});