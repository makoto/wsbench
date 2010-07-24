class Distribution
  def initialize
    @n, @sum_x, @sum_x_2 = 0, 0, 0
    @min, @max = nil, nil
  end
  
  def <<(x)
    @n += 1
    @sum_x += x
    @sum_x_2 += x**2
    @min = x if @min.nil? || x < @min
    @max = x if @max.nil? || x > @max
  end
  
  def stats
    mean_x = @sum_x.to_f / @n
    
    mean_x_2 = (1.0/@n) * @sum_x_2
    
    sd = Math.sqrt(mean_x_2 - mean_x**2)
    
    hash = {
      :count => @n, 
      :mean => mean_x,
      :sd => sd,
      :min => @min,
      :max => @max,
      :sum => @sum_x
    }
    hash.reduce({}){|s, c| s[c[0]] = ("%.3f" % (c[1])).to_f; s }
  end
  
  def self.stats(array)
    if array && array.size > 0
      dist = Distribution.new
      array.compact.each { |i| dist << i }
      dist.stats
    else
      nil
    end
  end
  
end

# Usage
# dist = Distribution.new
# 
# (0..10).each { |i| dist << i }
# 
# p dist.stats
