# wsbench - websocket benchmarking programme


    $ ruby wsbench.rb -h
    Usage: wsbench.rb [options] [ws[s]://]hostname[:port]/path 
     eg: wsbench.rb -c 10 ws://localhost:8080/echo
        -c connection                    number of concurrent connections (default: 10)
        -t timeout                       timeout (default: 30)
        -m message                       message size(bytes), (default: 1)
        -T type                          type echo|broadcast (default: echo)
        -h

TODO:
- Support for draft76 (currently 75 only)
- Support for wss (tls)
- Implement mullticast (= channels/rooms)
