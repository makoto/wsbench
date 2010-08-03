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

func echoServer(ws *websocket.Conn) { 
  // fmt.Printf("echoing ws: %v", ws.Protocol)
  io.Copy(ws, ws) 
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

func broadcastHandler(ws *websocket.Conn) {
    defer func() {
        subscriptionChan <- subscription{ws, false}
        ws.Close()
        fmt.Printf("Closed \n")
    }()

    fmt.Printf("Adding to subscription \n")

    subscriptionChan <- subscription{ws, true}
    fmt.Printf("Added to subscription \n")

    buf := make([]byte, 256)
    for {
        n, err := ws.Read(buf)
        fmt.Printf("Reding message %v \n", n)
        
        if err != nil {
            break
        }
        messageChan <- buf[0:n]
    }
}

func startServer() {
  http.Handle("/echo", websocket.Handler(echoServer))
  http.Handle("/broadcast", websocket.Handler(broadcastHandler))
  go http.ListenAndServe(":5555", nil)
}

func TestEcho(t *testing.T) {
  once.Do(startServer)
  msg := []byte("hello, world!")
  ws, err := websocket.Dial("ws://0.0.0.0:5555/echo", "", "http://0.0.0.0/")
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
  wsClients := &WSBench{Connections: 2, target: "ws://0.0.0.0:5555/echo"}
  if wsClients.Connections != 2 {
    t.Errorf("Setting connections ", wsClients.Connections, 2)
  }
  if wsClients.target != "ws://0.0.0.0:5555/echo" {
    t.Errorf("Setting target ", wsClients.target, "ws://0.0.0.0:5555/echo")
  }

}

func TestRunCreatesMultipleConnections(t *testing.T) {
  var ch = make(chan Result)

  once.Do(startServer)
  wsClients := &WSBench{Connections: 2, Ch: ch}
  wsClients.Run()
  fmt.Printf("A: %v ", wsClients.results)
  if len(wsClients.results) != 2 {
    t.Errorf("Running WSBench w 2 connections should return 2 results ", wsClients.results, 2)
  }
}

func TestStatShouldHaveStat(t *testing.T) {
  var ch = make(chan Result)

  once.Do(startServer)
  wsClients := &WSBench{Connections: 3, Ch: ch}
  wsClients.Run()
  fmt.Printf("B: %v ", wsClients.results)
  fmt.Printf("C: %v ", wsClients.Stats)
  if wsClients.Stats["sum"] <= 0 {
    t.Errorf("Stats should have sum", wsClients.Stats["sum"], 0)
  }

  if wsClients.Stats["count"] <= 0 {
    t.Errorf("Stats should have count", wsClients.Stats["count"], 0)
  }

  if wsClients.Stats["avg"] <= 0 {
    t.Errorf("Stats should have avg", wsClients.Stats["avg"], 0)
  }

  if wsClients.Stats["max"] <= 0 {
    t.Errorf("Stats should have max", wsClients.Stats["max"], 0)
  }

  if wsClients.Stats["min"] <= 0 {
    t.Errorf("Stats should have min", wsClients.Stats["min"], 0)
  }

}

func TestMoreThan30Connections(t *testing.T) {
  var ch = make(chan Result)

  once.Do(startServer)
  wsClients := &WSBench{Connections: 30, Ch: ch}
  wsClients.Run()
  // fmt.Printf("A: %v ", wsClients.results)
  if len(wsClients.results) < 30 {
    t.Errorf("Running WSBench w 2 connections should return 30 results ", len(wsClients.results), 30)
  }
}
