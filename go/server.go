package main


import (
        "http"
        "io"
        "websocket"
        "fmt"
)

// Echo the data received on the Web Socket.
// func EchoServer(ws *websocket.Conn) {
//         fmt.Printf("ws: %+v", ws)
//         io.Copy(ws, ws)
// }
func EchoServer(ws *websocket.Conn) { 
  // fmt.Printf("ws: %+v", ws)
  io.Copy(ws, ws) 
}



func main() {
        fmt.Printf("starting websocket server\n")
        // http.Handle("/echo", websocket.Handler(EchoServer));
        http.Handle("/echo", websocket.Draft75Handler(EchoServer));
        fmt.Printf("listening\n")
        err := http.ListenAndServe(":8080", nil);
        if err != nil {
                println("hello")
                panic("ListenAndServe: " + err.String())
        }
        fmt.Printf("end of main\n")
}

