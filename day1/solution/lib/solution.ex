defmodule Solution do

  def main(args) do
    input = "input.txt"

    modules = process_input(input)

    modules |> part1 |> IO.puts
    modules |> part2 |> IO.puts
  end

  def part1(modules) do
    modules |> Enum.map(&(floor(&1 / 3) - 2)) |> Enum.sum
  end

  def part2(modules) do
    modules |> Enum.map(&(fuel(&1))) |> Enum.sum
  end


  def fuel(mass) do
    mass
      |> Stream.iterate(&(floor(&1 / 3) - 2))
      |> Enum.take_while(&(&1 > 0))
      |> Enum.drop(1)
      |> Enum.sum
  end

  def process_input(path) do
    stream = File.stream!(path)

    stream |> Enum.map(&String.to_integer(String.trim(&1)))
  end

end
