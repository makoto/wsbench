require 'rubygems'
require 'em-websocket'
# Epoll will initially be configured to 1024 descriptors on ubuntu
desired_descriptors = 8192 * 2
file_descriptors = EventMachine.set_descriptor_table_size(desired_descriptors)
if file_descriptors == desired_descriptors
  p "Epoll configured with #{file_descriptors} file descriptors"
else
  p "Epoll configured with only #{file_descriptors} file descriptors"
end
port = (ARGV[0] && ARGV[0].to_i) || 8080
p port

EventMachine.epoll
EventMachine.run {

  EventMachine::WebSocket.start(:host => "0.0.0.0", :port => port) do |ws|
      ws.onopen {
        puts "WebSocket connection open"

        # publish message to the client
        # ws.send "Hello Client"
      }

      ws.onclose { puts "Connection closed" }
      ws.onmessage { |msg|
        puts "Recieved message: #{msg}"
        ws.send msg
      }
  end
}
