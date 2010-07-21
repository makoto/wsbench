/*
  The wsbench package implements ... bananas?
*/
package wsbench

import (
  "bytes"
  "fmt"
  "time"
  "websocket"
  "strconv"
)

func sum(a []int64) int64 { // returns an int
  var s int64 = 0
  for i := 0; i < len(a); i++ {
    s += a[i]
  }
  return s
}

func max(a []int64) int64 { // returns an int
  var s int64 = a[0]
  for i := 1; i < len(a); i++ {
    if a[i] >= s {
      s = a[i]
    }
  }
  return s
}

func min(a []int64) int64 { // returns an int
  var s int64 = a[0]
  for i := 1; i < len(a); i++ {
    if a[i] <= s {
      s = a[i]
    }
  }
  return s
}

type Result struct {
  time int64
}

type WSBench struct {
  connections int
  target      string
  results     []Result
  stats       map[string]int64
  ch          chan Result
}

// ch1 := make(chan int)
// var ch = make(chan Result)


func do_test(w *WSBench, msg []byte) {
  for i := 0; i < w.connections; i++ {
    start := time.Nanoseconds()
    ws, err := websocket.Dial("ws://localhost:12345/echo", "", "http://localhost/")
    if err != nil {
      panic("Dial failed: " + err.String())
    }

    if _, err := ws.Write([]byte(msg)); err != nil {
      panic("Write failed: " + err.String())
    }

    var response = make([]byte, 512)
    fmt.Println(bytes.Equal(response, response))
    n, err := ws.Read(response)

    if err != nil {
      panic("Read failed: " + err.String())
    }
    if !bytes.Equal(msg, response[0:n]) {
      panic("Message not the same: " + err.String())
    }
    ws.Close()
    delta := time.Nanoseconds() - start
    // Adding Dummy result for now.
    w.ch <- Result{time: delta}
  }
  close(w.ch)
}


func (w *WSBench) Run() {
  w.results = make([]Result, w.connections)
  msg := []byte("hello, world!")

  go do_test(w, msg)

  for i := 0; i < w.connections; i++ {
    m := <-w.ch
    fmt.Printf("i: %+v m: %+v", i, m)
    w.results[i] = m
    if closed(w.ch) {
      fmt.Println("Finished 2\n")
      break
    }
  }

  times := make([]int64, w.connections)
  for i := range w.results {
    fmt.Printf("i: %v time: %v\n", i, w.results[i].time)
    times[i] = w.results[i].time
  }
  lenS := strconv.Itoa(len(times))
  len64, _ := strconv.Btoi64(lenS, 10)

  var avg int64 = sum(times) / len64

  w.stats = map[string]int64{
    "sum":   sum(times),
    "avg":   avg,
    "max":   max(times),
    "min":   min(times),
    "count": len64}
}