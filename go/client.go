package main

import (
	"websocket"
	"fmt"
  // "strings"
)

func main() {
  println(1)
 	ws, err := websocket.Dial("ws://localhost:8080/", "", "http://localhost:8080/");
 	if err != nil {
		panic("Dial: " + err.String())
	}
	println(2)
	if _, err := ws.Write([]byte("hello, world!\n")); err != nil {
		panic("Write: " + err.String())
	}
	var msg = make([]byte, 512);
	println(3)
	if output, err := ws.Read(msg); err != nil {
		panic("Read: " + err.String())
	}else{
  	println(4)
  	fmt.Printf("msg: %s\n", msg)
  	fmt.Printf("output: %v\n", output)
	}
	println(5)
	// use msg[0:n]
}
