package main

import (
  // "bytes"
  "fmt"
  "http"
  "io"
  // "log"
  // "net"
  "wsbench"
  "websocket"
)
func echoServer(ws *websocket.Conn) { io.Copy(ws, ws) }

func startServer() {
  http.Handle("/echo", websocket.Handler(echoServer))
  go http.ListenAndServe(":12345", nil)
}

func main() {
  go startServer()
  var ch = make(chan wsbench.Result)
  
  ws := &wsbench.WSBench{Connections:2, Ch: ch}
  ws.Run()
  
  fmt.Printf("A: %v :B", ws.Stats)
}

