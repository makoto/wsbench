package wsbench

import (
  "bytes"
  "fmt"
  "http"
  "io"
  // "log"
  // "net"
  "once"
  "websocket"
  "testing"
)

var serverAddr string

func echoServer(ws *websocket.Conn) { io.Copy(ws, ws) }

func startServer() {
  http.Handle("/echo", websocket.Handler(echoServer))
  go http.ListenAndServe(":12345", nil)
}

func TestEcho(t *testing.T) {
  once.Do(startServer)
  msg := []byte("hello, world!")
  ws, err := websocket.Dial("ws://localhost:12345/echo", "", "http://localhost/")
  if err != nil {
    t.Errorf("WebSocket handshake: %v", err)
  }
  if _, err := ws.Write([]byte(msg)); err != nil {
    t.Errorf("Write: error %v", err)
  }
  var response = make([]byte, 512)
  fmt.Println(bytes.Equal(response, response))
  n, err := ws.Read(response)

  if err != nil {
    t.Errorf("Read: error %v", err)
  }

  if !bytes.Equal(msg, response[0:n]) {
    t.Errorf("Echo: expected %q got %q", msg, response)
  }
  ws.Close()
}

func TestSetClientConnections(t *testing.T) {
  once.Do(startServer)
  wsClients := &WSBench{connections: 2, target: "ws://localhost:12345/echo"}
  if wsClients.connections != 2 {
    t.Errorf("Setting connections error")
  }
  if wsClients.target != "ws://localhost:12345/echo" {
    t.Errorf("Setting connections error")
  }

}

// NOTE: Looks I can do only one channel?

func TestRunCreatesMultipleConnections(t *testing.T) {
  var ch = make(chan Result)

  once.Do(startServer)
  wsClients := &WSBench{connections: 2, ch: ch}
  wsClients.Run()
  fmt.Printf("A: %v ", wsClients.results)
  if len(wsClients.results) != 2 {
    t.Errorf("Running WSBench w 2 connections should return 2 results ", wsClients.results, 2)
  }
}

func TestStatShouldHaveSum(t *testing.T) {
  var ch = make(chan Result)

  once.Do(startServer)
  wsClients2 := &WSBench{connections: 2, ch: ch}
  wsClients2.Run()
  fmt.Printf("B: %v ", wsClients2.results)
  fmt.Printf("C: %v ", wsClients2.stats["sum"])
  if wsClients2.stats["sum"] <= 0 {
    t.Errorf("Stats should have sum", wsClients2.stats["sum"], 0)
  }
}
