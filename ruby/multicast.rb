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

class MyChannel
  attr_reader :ws, :channel_id
  def initialize(opts)
    @ws = opts[:ws]
    @channel_id = opts[:channel_id]
  end
end

@channels = []
@channel = nil
@sid = nil
c = 0
EventMachine.epoll
EventMachine.run {
  EventMachine::WebSocket.start(:host => "0.0.0.0", :port => port) do |ws|
    ws.onopen {
      c = c + 1
      p "channel: #{c}"
      # @channels << MyChannel.new(:ws => ws, :channel_id => c)
      p 1
      @channel = EM::Channel.new
      p 2
      @channels << @channel
      p 3
      @sid = @channel.subscribe { |msg| ws.send msg }
      p "CHANNEL: #{@channels.size}"
    }    
    
    ws.onmessage { |msg|
      p "onmessage: #{@channels.size} channels"
      # @channels.each do |c|
      #   p "cid: #{c.channel_id}"
      #   c.ws.send msg
      # end
      # ws.send msg
      # p @channel.channel_id
      @channels.each {|c|
        c.push msg
      }
      # @channel.push msg      
    }

    ws.onclose {
      p "onclose: #{@channel} #{@sid}"
      # @channel.unsubscribe(@sid)
    }
    
  end

  puts "Server started"
}
