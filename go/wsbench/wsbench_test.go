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
	print(msg)

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
