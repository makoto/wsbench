#!/usr/bin/ruby

require "rubygems"
require 'eventmachine'
require 'em-http'
require 'uri'
require 'json'

uri = URI.parse('ws://localhost')
uri.port = 8080
uri.path = "/"
start_time = Time.now
results = []
EM.run {
  200.times do |i|
    ws = EventMachine::HttpRequest.new(uri.to_s).get(:timeout => 10)
    ws.callback{
      results << (Time.now - start_time)
    }
    ws.stream{}
    
  end
  EventMachine::add_timer(10) {
    #puts "Failed to recevie message from Pusher -- FAIL"
    max = results.max
    min = results.min
    avg = results.reduce{|total, current| total = total + current; total} / results.size
    p "#{results.size } connections, max #{max}, min #{min}, avg #{avg}"
    EventMachine::stop
  }
}