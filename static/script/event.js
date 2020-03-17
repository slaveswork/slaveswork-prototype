
window.onload = function(){
    var start = document.getElementById("electron_start");
    console.log(start), start.addEventListener("astilectron-ready", (function () {
        astilectron.sendMessage("hello", (function (e) {
            console.log("received " + e)
        }))
    }))
};
