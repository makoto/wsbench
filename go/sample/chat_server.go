//  http://gist.github.com/505923
package main

import (
    "flag"
    "http"
    "log"
    "template"
    "websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
    flag.Parse()
    go hub()
    http.HandleFunc("/", homeHandler)
    http.Handle("/ws", websocket.Draft75Handler(clientHandler))
    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Exit("ListenAndServe:", err)
    }
}

var messageChan = make(chan []byte)

type subscription struct {
    conn      *websocket.Conn
    subscribe bool
}

var subscriptionChan = make(chan subscription)

func hub() {
    conns := make(map[*websocket.Conn]int)
    for {
        select {
        case subscription := <-subscriptionChan:
            conns[subscription.conn] = 0, subscription.subscribe
        case message := <-messageChan:
            for conn, _ := range conns {
                if _, err := conn.Write(message); err != nil {
                    conn.Close()
                }
            }
        }
    }
}

func clientHandler(ws *websocket.Conn) {
    defer func() {
        subscriptionChan <- subscription{ws, false}
        ws.Close()
    }()

    subscriptionChan <- subscription{ws, true}

    buf := make([]byte, 256)
    for {
        n, err := ws.Read(buf)
        if err != nil {
            break
        }
        messageChan <- buf[0:n]
    }
}

// Handle home page requests.
func homeHandler(c *http.Conn, req *http.Request) {
    homeTempl.Execute(req.Host, c)
}

var homeTempl *template.Template

func init() {
    homeTempl = template.New(nil)
    homeTempl.SetDelims("«", "»")
    if err := homeTempl.Parse(homeStr); err != nil {
        panic("template error: " + err.String())
    }
}

const homeStr = `
<html>
<head>
<title>Chat Example</title>
<script type="text/javascript" src="http://ajax.googleapis.com/ajax/libs/jquery/1.4.2/jquery.min.js"></script>
<script type="text/javascript">
    $(function() {

    var conn;
    var msg = $("#msg");
    var log = $("#log");

    function appendLog(msg) {
        var d = log[0]
        var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
        msg.appendTo(log)
        if (doScroll) {
            d.scrollTop = d.scrollHeight - d.clientHeight;
        }
    }

    $("#form").submit(function() {
        if (!conn) {
            return false;
        }
        if (!msg.val()) {
            return false;
        }
        conn.send(msg.val());
        msg.val("");
        return false
    });

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://«@»/ws");
        conn.onclose = function(evt) {
            appendLog($("<div><b>Connection closed.</b></div>"))
        }
        conn.onmessage = function(evt) {
            appendLog($("<div/>").text(evt.data))
        }
    } else {
        appendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
    }
    });
</script>
<style type="text/css">
html {
    overflow: hidden;
}

body {
    overflow: hidden;
    padding: 0;
    margin: 0;
    width: 100%;
    height: 100%;
    background: gray;
}

#log {
    background: white;
    margin: 0;
    padding: 0.5em 0.5em 0.5em 0.5em;
    position: absolute;
    top: 0.5em;
    left: 0.5em;
    right: 0.5em;
    bottom: 3em;
    overflow: auto;
}

#form {
    padding: 0 0.5em 0 0.5em;
    margin: 0;
    position: absolute;
    bottom: 1em;
    left: 0px;
    width: 100%;
    overflow: hidden;
}

</style>
</head>
<body>
<div id="log"</div>
<form id="form">
    <input type="submit" value="Send" />
    <input type="text" id="msg" size="64"/>
</form>
</body>
</html> `