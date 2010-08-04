package main

import (
	"http"
	"io"
	"websocket"
	"fmt"
)

// Echo the data received on the Web Socket.
func EchoServer(ws *websocket.Conn) {
  fmt.Printf("at EchoServer")
	io.Copy(ws, ws);
}

func main() {
  http.Handle("/echo", websocket.Handler(EchoServer));
	err := http.ListenAndServe(":8080", nil);
	if err != nil {
		panic("ListenAndServe: " + err.String())
	}
}
