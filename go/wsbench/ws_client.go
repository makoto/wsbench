package main

import (
  "fmt"
  "http"
  "io"
  "wsbench"
  "websocket"
  "flag"
)
func echoServer(ws *websocket.Conn) { io.Copy(ws, ws) }

func startServer() {
  http.Handle("/echo", websocket.Handler(echoServer))
  go http.ListenAndServe(":12345", nil)
}

func main() {
  go startServer()
  var ch = make(chan wsbench.Result)

  var c *int = flag.Int("c", 1, "number of concurrent connections")

  flag.Parse()
  fmt.Println("c has value ", *c);
  
  for i := 0; i < flag.NArg(); i++ {
    fmt.Printf("Flag: %v :", flag.Arg(i))
  }
  
  
  ws := &wsbench.WSBench{Connections:*c, Ch: ch}
  ws.Run()
  
  fmt.Printf("A: %v :B", ws.Stats)
}

