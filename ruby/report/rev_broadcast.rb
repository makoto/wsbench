require 'rubygems'
require 'rev/websocket'

$connections = []
@sid = nil


class MyConnection < Rev::WebSocket
  def on_open
    puts "WebSocket opened"
    $connections << self
  end

  def on_message(data)
#    puts "WebSocket data received: '#{data}'"
    $connections.each do |c|
      c.send_message data
    end
  #  send_message data
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
