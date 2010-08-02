package main

import (
	"fmt"
	"websocket"
)

const message = "A message"

func main() {
	ws, err := websocket.Dial("ws://127.0.0.1:12345/match", "", "")
	if err != nil {
		panic(err)
	}
	if _, err := ws.Write([]byte(message)); err != nil {
		panic(err)
	}
	var resp = make([]byte, 512)
	n, err := ws.Read(resp)
	if err != nil {
		panic(err)
	}
	fmt.Println("Received:", string(resp[0:n]))
}
