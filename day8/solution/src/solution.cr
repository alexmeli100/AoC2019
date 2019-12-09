H = 6
W = 25

def read_input(path)
  input = File.read(path)

  input.chars()
  
end

def get_layers(pixels, size) 
  pixels.in_groups_of(size)
end

def fewest(layers) 
  min = layers.min_by { |l| l.count('0')}

  min.count('1') * min.count('2')
end

def get_picture(layers) 
  picture = layers.last

  layers.reverse.each { |layer| 
    (0...layer.size).each { |i| 
      next if layer[i] == '2'
      picture[i] = layer[i]
    }
  }

  picture
end

def print_picture(picture) 
  get_layers(picture, W).each { |p| 
    puts p.join.gsub("0", " ").gsub("1", "#")
  }
end

input = read_input("input.txt")
layers = get_layers(input, H*W)
puts fewest(layers)
print_picture(get_picture(layers))