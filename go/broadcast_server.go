// Majority of the code is from 
// http://gary.beagledreams.com/page/go-websocket-chat
package main


import (
    "flag"
    "http"
    "log"
    "websocket"
    "fmt"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
    flag.Parse()
    go hub()
    http.Handle("/broadcast", websocket.Handler(clientHandler))
    if err := http.ListenAndServe(*addr, nil); err != nil {
        log.Exit("ListenAndServe:", err)
    }
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

func clientHandler(ws *websocket.Conn) {
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
