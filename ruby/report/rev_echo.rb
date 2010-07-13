require 'rubygems'
require 'rev/websocket'

class MyConnection < Rev::WebSocket
  def on_open
    puts "WebSocket opened"
#    send_message("Hello, world!")
  end

  def on_message(data)
    puts "WebSocket data received: '#{data}'"
    send_message(data)
  end

  def on_close
    puts "WebSocket closed"
  end
end

host = '0.0.0.0'
port = ARGV[0]
p port

server = Rev::WebSocketServer.new(host, port, MyConnection)
server.attach(Rev::Loop.default)

Rev::Loop.default.run
