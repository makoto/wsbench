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
p connections

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

uri = URI.parse('ws://localhost')
uri.port = 8080
uri.path = "/"
results = []
begin_time = Time.now.to_f
EM.run {
  connections.times do |i|
    
    start_time = Time.now.to_f
    ws = EventMachine::HttpRequest.new(uri.to_s).get(:timeout => 10)
    ws.callback{
      end_time = Time.now.to_f
      result = Connection.new(:connection_id => i, :start_time => start_time, :end_time => end_time)
      results << result
    }
    ws.stream{
      
    }
    
  end
  EventMachine::add_timer(10) {
    # p results
    oredered_results = results.sort{|a, b| a.time_taken <=> b.time_taken}
    max = oredered_results.last.time_taken
    min = oredered_results.first.time_taken
    sum = results.reduce(0){|total, current| total = total + current.time_taken; total}
    avg = sum / results.size
    p "#{results.size } connections, sum #{sum}, max #{max}, min #{min}, avg #{avg}"
    EventMachine::stop
  }
}