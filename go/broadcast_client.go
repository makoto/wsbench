package main

import (
	"fmt"
	"websocket"
	list "container/list"
)

const message = "A message"

func sub(ws *websocket.Conn)  {
  println("Beginning Sub")
  
  var resp = make([]byte, 512)
  for {
    println("Waiting Read")
    _, err := ws.Read(resp)

    if err != nil {
     panic(err)
    }
    // fmt.Println("Received 1:", string(resp[0:n]))
    // fmt.Println("Received %v:", n)
    ch <- 1
  }
}

func pub()  {
  ws, err := websocket.Dial("ws://127.0.0.1:8080/broadcast", "", "")
  if err != nil {
   panic(err)
  }
  println("About to write")
  if _, err := ws.Write([]byte(message)); err != nil {
   panic(err)
  }
}

var ch = make(chan int)

func main() {
  ws1, err := websocket.Dial("ws://127.0.0.1:8080/broadcast", "", "")
  if err != nil {
   panic(err)
  }
  go sub(ws1)

  ws2, err := websocket.Dial("ws://127.0.0.1:8080/broadcast", "", "")
  if err != nil {
   panic(err)
  }
  go sub(ws2)

  ws3, err := websocket.Dial("ws://127.0.0.1:8080/broadcast", "", "")
  if err != nil {
   panic(err)
  }
  go sub(ws3)


  pub()
  l := list.New()
  for {
    res := <- ch
    
    l.PushBack(1)
  	fmt.Printf("%v",l.Len())
    fmt.Println("Received: %v,  %+v", l.Len(), res)
  }
}
