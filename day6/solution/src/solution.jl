#=
solution:
- Julia version: 1.2.0
- Author: Alex
- Date: 2019-12-06
=#

using DataStructures

mutable struct Orbit
    children::Dict
    parents::Dict

    function Orbit(planets)
        c = Dict{String, Array{String}}()
        p = Dict{String, String}()

        for s in planets
            p[s[2]] = s[1]
            c[s[1]] = push!(get(c, s[1], []), s[2])
        end

        new(c, p)
    end
end

function path_to(orbit::Orbit, src::String, dst::String)
    path = []
    current = src
    while haskey(orbit.parents, current)
        push!(path, orbit.parents[current])
        current = orbit.parents[current]
    end

    path
end

function get_distances(orbit::Orbit, root::String)
    q = Queue{String}()
    distances = Dict{String, Int}(root => 0)

    for c in orbit.children[root]
        enqueue!(q, c)
    end

    while !isempty(q)
        x = dequeue!(q)
        distances[x] = distances[orbit.parents[x]] + 1

        for c in get(orbit.children, x, [])
            enqueue!(q, c)
        end
    end

    distances
end

function read_input()
    f = open("input.txt", "r")
    input = read(f, String)

    map((s) -> split(s, ")"), split(input, "\n"))
end

function solve()
    planets = read_input()
    o = Orbit(planets)
    distances = get_distances(o, "COM")

    you = path_to(o, "YOU", "COM")
    san = path_to(o, "SAN", "COM")
    lca = intersect(you, san)[1]

    res = (distances["YOU"] + distances["SAN"]) - 2 * distances[lca] - 2
    println(sum(values(distances)))
    println(res)
end

solve()