/*
  The wsbench package implements ... bananas?
*/
package wsbench

type Result struct{
  time int
}

type WSBench struct{
  connections int
  results []Result
}

func (w *WSBench) Run (){
  w.results = make([]Result, w.connections)
  for i := 0; i < 2; i++ {
    // Adding Dummy result for now.
    w.results[i] = Result{time:2}
  }
}


