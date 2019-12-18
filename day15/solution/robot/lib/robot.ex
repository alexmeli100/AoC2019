
defmodule Point do
  defstruct X: 0, Y: 0

  def neighbor(%Point{X: x, Y: y}, d) do
    case d do
      1 -> %Point{X: x, Y: y-1}
      2 -> %Point{X: x, Y: y+1}
      3 -> %Point{X: x-1, Y: y}
      4 -> %Point{X: x+1, Y: y}
    end
  end
end

defmodule Robot do

  @reverse %{1 => 2, 2 => 1, 3 => 4, 4 => 3}


  def get_map(sock) do
    m = %{}
    trace = []
    origin = %Point{}
    oxygen = nil

    explore(sock, m, trace, oxygen, origin)
  end

  defp explore(sock, m, trace, oxygen, p) do
    {dir, t} = find_dir(m, p, trace)

    if !dir do
      :gen_tcp.close(sock)
      {oxygen, m}
    else
      next = Point.neighbor(p, dir)
      input = get_input(sock, dir)

      next_trace = if input > 0 && !Map.has_key?(m, next), do: [@reverse[dir] | t], else: t
      oxy_pos = if input > 1, do: next, else: oxygen
      next_pos = if input > 0, do: next, else: p
      explored = Map.put(m, next, input)

      explore(sock, explored, next_trace, oxy_pos, next_pos)
    end
  end

  defp find_dir(m, p, trace) do
    dir = find_unexplored(m, p)

    case !dir do
      true -> List.pop_at(trace, 0)
      _ -> {dir, trace}
    end
  end

  def get_input(sock, dir) do
    case :gen_tcp.send(sock, "#{dir}\n") do
      :ok -> :ok
      {:error, reason} ->
        IO.puts "Failed to write output #{inspect reason}"
        Process.exit(self(), reason)
    end


    case :gen_tcp.recv(sock, 0) do
      {:ok, res} ->
         String.to_integer(String.trim(res))

      {:error, reason} ->
        IO.puts "Failed to read input #{inspect reason}"
        Process.exit(self(), reason)
    end
  end

  def open_connection() do
    case :gen_tcp.connect('localhost', 8080, [:binary,  packet: :line, active: false]) do
      {:ok, sock} ->
        sock

      {:error, reason} ->
        IO.puts "Failed to open connection #{inspect reason}"
        Process.exit(self(), reason)
    end
  end

  defp find_unexplored(m, p) do
    1..4
    |> Enum.filter(&(!Map.has_key?(m, Point.neighbor(p, &1))))
    |> Enum.at(0)
  end

  def run() do
    sock = open_connection()

    {oxygen, m} = get_map(sock)
    dis = flood_fill(m, %Point{})
    oxy_fill = flood_fill(m, oxygen)
    max = oxy_fill |> Map.values |> Enum.max

    IO.puts dis[oxygen]
    IO.puts max
  end

  def flood_fill(m, origin) do
    q = Deque.new(200) |> Deque.append(origin)
    dis = %{origin => 0}


    flood_fill(m, dis, q)
  end

  def flood_fill(m, dis, q) do
    {front, q1} = Deque.popleft(q)

    if !front do
      dis
    else
      neighbors = 1..4
        |> Enum.map(&(Point.neighbor(front, &1)))
        |> Enum.filter(&(!Map.has_key?(dis, &1) && Map.get(m, &1, 0) > 0))

      d1 = Map.get(dis, front) + 1

      {next_dis, next_q} = neighbors
        |> Enum.reduce({dis, q1}, fn n, {d, q} -> {Map.put(d, n, d1), Deque.append(q, n)} end)

      flood_fill(m, next_dis, next_q)
    end
  end
end

Robot.run()




