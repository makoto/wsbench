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

RESULT:
- The full result is at http://github.com/makoto/wsbench/tree/master/ruby/report/
- The summary of comparing node(broadcast-server.js) , rev-websocket(rev_braodcast.rb) and em-websocket (multicast.rb) is at http://bit.ly/d4W2hV
- This is run on 24th July 2010 and 8th Aug
- Run on 2 amazon ec2 (small) instances , servers in one machine, and wsbench.rb in another.