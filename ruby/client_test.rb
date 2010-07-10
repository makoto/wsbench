#!/usr/bin/ruby

#  (1..25).to_a.map{|a| a * 10}.each do | i|
#  p  `ruby ./wsbench.rb #{i}`
#  end



require "rubygems"
require 'eventmachine'
require 'em-http'
require 'uri'
p ARGV

uri = URI.parse(ARGV.first)

EM.epoll
EM.run {
  start_time = Time.now.to_f
  ws = EventMachine::HttpRequest.new(uri.to_s).get(:timeout => 10)
  ws.callback{
    p "connected"
    ws.send("hello")
  }
  ws.stream{|msg|
    p "received"
    p msg
  }
}
