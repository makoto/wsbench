package main

import (
  "fmt"
  "http"
  "io"
  "once"
  "websocket"
)

var serverAddr string

func echoServer(ws *websocket.Conn) { 
  io.Copy(ws, ws) 
}

func startServer() {
  http.Handle("/echo", websocket.Handler(echoServer))
  go http.ListenAndServe(":5555", nil)
}

func main() {
  once.Do(startServer)
  for i := 0; i < 30; i++ {
    fmt.Println(i)
    // body
    _, err := websocket.Dial("ws://localhost:5555/echo", "", "http://localhost/")
    if err != nil {
      panic("Dial failed: " + err.String())
    }
  }
}

