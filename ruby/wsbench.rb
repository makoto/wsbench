#!/usr/bin/ruby

#  (1..25).to_a.map{|a| a * 10}.each do | i|
#  p  `ruby ./wsbench.rb #{i}`
#  end



require "rubygems"
require 'eventmachine'
require 'em-http'
require 'uri'
require 'json'
require 'ostruct'
require 'optparse'

options = OpenStruct.new
options.connections = 10
options.timeout = 30
options.message = 1
options.descriptors = 8192

OptionParser.new do |o|
  o.banner = "Usage: wsbench.rb [options] [ws[s]://]hostname[:port]/path \n eg: wsbench.rb -c 10 ws://localhost:8080/echo"
  o.on('-c connection ', 
    help = "number of concurrent connections (default: #{options.connections})") {|s|options.connections = s.to_i}
  o.on('-t timeout', help = "timeout (default: #{options.timeout})"){|s| options.timeout = s.to_f}
  o.on('-d maxdescrtors', help = "max number of descriptors (default: #{options.descriptors})"){|s| options.descriptors = s.to_i}
  o.on('-m mesage', help ="message size(bytes), (default: #{options.message})"){|s| options.message = s.to_f}
  o.on('-h') { puts o; exit }
  o.parse!
end

p ARGV
p options

uri = URI.parse(ARGV.first)

connections = options.connections
timeout = options.timeout
message = options.message
descriptors = options.descriptors

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

results = []
results2 = []

desired_descriptors = descriptors
file_descriptors = EventMachine.set_descriptor_table_size(desired_descriptors)
if file_descriptors == desired_descriptors
  p "Epoll configured with #{file_descriptors} file descriptors"
else
  p "Epoll configured with only #{file_descriptors} file descriptors"
end

begin_time = Time.now.to_f
EM.epoll
EM.run {
  connections.times do |i|
    
    start_time = Time.now.to_f
    ws = EventMachine::HttpRequest.new(uri.to_s).get(:timeout => 10)
    ws.callback{
      end_time = Time.now.to_f
      result = Connection.new(:connection_id => i, :start_time => start_time, :end_time => end_time)
      results << result
      ws.send(JSON.generate({:connection_id => i, :start_time => end_time, :data => "a" * message }))
    }
    ws.stream{|msg|
      reply = JSON.parse(msg)
      result = Connection.new(
        :connection_id => reply["connection_id"], 
        :start_time => reply["start_time"],
        :end_time => Time.now.to_f
      )
      results2 << result
    }
  end

  EventMachine::add_periodic_timer(1) {
  if results2.size == connections
    print "SUCCESS (#{results2.size}/#{connections}), "
    print "Connect: #{show_result(results)} ,"
    print "Message: #{show_result(results2)}"

    EventMachine::stop
  end
  }
  
  EventMachine::add_timer(timeout) {
    print "TIMEOUT (#{results2.size} / #{connections}), "
    print "Connect: #{show_result(results)} ,"
    print "Message: #{show_result(results2)}"

    EventMachine::stop
  }
}
end_time = Time.now.to_f
print ", Total #{sprintf("%.3f", end_time - begin_time)}"
