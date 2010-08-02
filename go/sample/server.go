package main

import (
	"http"
	"io"
	"fmt"
	"websocket"
)

type joinRequest struct {
	ws  *websocket.Conn      // the senders Conn
	ret chan *websocket.Conn // a chan on which to receive a partner's Conn
}

var requests = make(chan *joinRequest)

// handler accepts a websocket connection, and requests a partner
// it then copies everything from the partner's conn to its own conn
func handler(ws *websocket.Conn) {
  fmt.Println("Handling:")
	
	r := make(chan *websocket.Conn)
	requests <- &joinRequest{ws, r}
	p_ws := <-r
	io.Copy(ws, p_ws)
}

// matchmaker exchanges conns between sequential pairs of joinRequests
func matchmaker() {
	var p *joinRequest // = nil
	for r := range requests {
		if p == nil {
			// we have no pending requests, 
			// so this one becomes pending
			p = r
			continue // wait for the next joinRequest
		}
		// p != nil, so let's match p with r
		p.ret <- r.ws
		r.ret <- p.ws
		p = nil // reset p, as we have no pending requests now
	}
}

func main() {
	go matchmaker()
	http.Handle("/match", websocket.Handler(handler))
	http.ListenAndServe(":12345", nil) // this call will block forever
}

