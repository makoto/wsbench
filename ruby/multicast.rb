require 'rubygems'
require 'em-websocket'
# Epoll will initially be configured to 1024 descriptors on ubuntu
desired_descriptors = 8192 * 4
file_descriptors = EventMachine.set_descriptor_table_size(desired_descriptors)
if file_descriptors == desired_descriptors
  p "Epoll configured with #{file_descriptors} file descriptors"
else
  p "Epoll configured with only #{file_descriptors} file descriptors"
end
port = (ARGV[0] && ARGV[0].to_i) || 8080
p port

EventMachine.run {
  @channels = {}

  EventMachine::WebSocket.start(:host => "0.0.0.0", :port => port, :debug => true) do |ws|
    @channel = nil
    @sid = nil
    ws.onopen {
      # p ws.request["Path"]
      channel_name = ws.request["Path"]
      # p @channels[channel_name]
      @channel = if @channels[channel_name]
        @channels[channel_name]
      else
        @channels[channel_name] = EM::Channel.new
      end
      @sid = @channel.subscribe { |msg| ws.send msg }
      p "subscribed #{@sid} to #{channel_name}"
    }    
    
    ws.onmessage { |msg|
      p "onmessage"
      @channel.push msg
    }

    ws.onclose {
      p "onclose: #{@channel} #{@sid}"
      # if @channel
      #   @channel.unsubscribe(@sid)
      # end
    }
    
  end

  puts "Server started"
}
