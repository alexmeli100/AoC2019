package main

import (
	"github.com/deckarep/golang-set"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Path struct {
	dir  rune
	step int
}

type Point struct {
	x, y int
}

type Wire struct {
	paths     []Path
	pos       Point
	grid      mapset.Set
	stepCount map[Point]int
}

var disMap = map[rune]Point{'U': {0, 1}, 'D': {0, -1}, 'L': {-1, 0}, 'R': {1, 0}}

func (w *Wire) drawPath() {
	length := 0

	for _, p := range w.paths {
		for i := 0; i < p.step; i++ {
			w.pos = w.pos.Add(disMap[p.dir])
			w.grid.Add(w.pos)
			length++

			if _, ok := w.stepCount[w.pos]; !ok {
				w.stepCount[w.pos] = length
			}
		}
	}
}

func intersections(w1 *Wire, w2 *Wire) []Point {
	p := w1.grid.Intersect(w2.grid).ToSlice()
	points := make([]Point, len(p))
	for i, v := range p {
		points[i] = v.(Point)
	}

	return points
}

func findCrossing(w1 *Wire, w2 *Wire) (Point, Point) {
	points := intersections(w1, w2)

	minDis := points[0]
	minStep := points[0]

	for _, p := range points[1:] {
		if p.mahnattanDis() < minDis.mahnattanDis() {
			minDis = p
		}

		c1 := w1.stepCount[p] + w2.stepCount[p]
		c2 := w1.stepCount[minStep] + w2.stepCount[minStep]

		if c1 < c2 {
			minStep = p
		}
	}

	return minDis, minStep
}

func (p Point) Add(other Point) Point {
	return Point{p.x + other.x, p.y + other.y}
}

func (p Point) mahnattanDis() int {
	return abs(p.x) + abs(p.y)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	wires := parseFile("input.txt")

	w1 := wires[0]
	w2 := wires[1]
	w1.drawPath()
	w2.drawPath()

	part1, part2 := findCrossing(&w1, &w2)
	log.Println(part1.mahnattanDis())
	log.Println(w1.stepCount[part2] + w2.stepCount[part2])
}

func parseFile(path string) []Wire {
	var wires []Wire

	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	input := strings.Fields(string(bytes))

	for _, s := range input {
		wires = append(wires, newWire(s))
	}

	return wires
}

func newWire(s string) Wire {
	var paths []Path
	p := strings.Split(s, ",")

	for _, x := range p {
		dir := rune(x[0])
		step, err := strconv.Atoi(string(x[1:]))

		if err != nil {
			log.Fatal(err)
		}

		paths = append(paths, Path{dir, step})
	}

	return Wire{paths, Point{0, 0}, mapset.NewSet(), map[Point]int{}}
}
