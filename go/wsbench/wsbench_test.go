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
  http.Handle("/echo", websocket.Handler(echoServer));
  go http.ListenAndServe(":12345", nil);
}

func TestEcho(t *testing.T) {
  once.Do(startServer)
	msg := []byte("hello, world!")
  ws, err := websocket.Dial("ws://localhost:12345/echo", "", "http://localhost/");
  if err != nil {
   t.Errorf("WebSocket handshake: %v", err)
  }
  if _, err := ws.Write([]byte(msg)); err != nil {
   t.Errorf("Write: error %v", err)
  }
  var response = make([]byte, 512)
  fmt.Println(bytes.Equal(response, response))
  n, err := ws.Read(response); 
  
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
 wsClients := &WSBench{connections:2}
 if wsClients.connections != 2{
  t.Errorf("Setting connections error")
 }
}

func TestRunCreatesMultipleConnections(t *testing.T) {
  once.Do(startServer)
  wsClients := &WSBench{connections:2}
  wsClients.Run()
  fmt.Printf("A: %v :B", wsClients.results)
  if len(wsClients.results) != 2 {
    t.Errorf("Running WSBench w 2 connections should return 2 results ", wsClients.results, 2)
  }
}


