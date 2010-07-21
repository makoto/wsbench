package main

import (
	"http"
	"io"
	"websocket"
)

// Echo the data received on the Web Socket.
func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws);
}

func main() {
	http.Handle("/echo", websocket.Handler(EchoServer));
	err := http.ListenAndServe(":12345", nil);
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}