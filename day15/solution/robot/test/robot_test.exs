defmodule RobotTest do
  use ExUnit.Case
  doctest Robot

  test "greets the world" do
    assert Robot.hello() == :world
  end
end
