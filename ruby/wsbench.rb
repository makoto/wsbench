#!/usr/bin/ruby

#  (1..25).to_a.map{|a| a * 10}.each do | i|
#  p  `ruby ./wsbench.rb #{i}`
#  end



require "rubygems"
require 'eventmachine'
require 'em-http'
require 'uri'
require 'json'

connections =  (ARGV[0] && ARGV[0].to_i) || 250

class Connection
  attr_accessor :start_time, :end_time
  def initialize(options)
    @start_time = options[:start_time]
    @end_time = options[:end_time]
  end
  def time_taken
    @end_time - @start_time
  end
end

def show_result(array)
  oredered_array = array.sort{|a, b| a.time_taken <=> b.time_taken}
  max = oredered_array.last.time_taken
  min = oredered_array.first.time_taken
  sum = array.reduce(0){|total, current| total = total + current.time_taken; total}
  avg = sum / array.size
  "sum #{sprintf("%.3f", sum)}, max #{sprintf("%.3f",max)}, min #{sprintf("%.3f",min)}, avg #{sprintf("%.3f", avg)}"
end


uri = URI.parse('ws://localhost')
uri.port = 8080
uri.path = "/"
results = []
results2 = []

begin_time = Time.now.to_f
EM.run {
  connections.times do |i|
    
    start_time = Time.now.to_f
    ws = EventMachine::HttpRequest.new(uri.to_s).get(:timeout => 10)
    ws.callback{
      end_time = Time.now.to_f
      result = Connection.new(:connection_id => i, :start_time => start_time, :end_time => end_time)
      results << result
      # ws.send(JSON.generate({:connection_id => i, :start_time => end_time, :data => "a" * 1000 }))
      ws.send(JSON.generate({:connection_id => i, :start_time => end_time}))
    }
    ws.stream{|msg|
      reply = JSON.parse(msg)
      # p reply
      result = Connection.new(
        :connection_id => reply["connection_id"], 
        :start_time => reply["start_time"],
        :end_time => Time.now.to_f
      )
      results2 << result
    }
    
  end
  EventMachine::add_timer(20) {
    print "#{results.size}: "
    print "Connect: #{show_result(results)} :"
    print "Message: #{show_result(results2)}\n"

    EventMachine::stop
  }
}

