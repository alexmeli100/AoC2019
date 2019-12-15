package main

import (
	"github.com/alexmeli100/AoC2019/day7/solution/intcode"
	"github.com/alexmeli100/AoC2019/goutils"
	"log"
	"sync"
)

const (
	Empty = iota
	Wall
	Block
	Paddle
	Ball
)

type Point struct {
	x int
	y int
}

type GameHandler struct {
	In           chan int
	Out          chan int
	Screen       map[Point]int
	BallCoords   Point
	PaddleCoords Point
}

func (h GameHandler) CountBlocks() int {
	count := 0

	for _, v := range h.Screen {
		if v == 2 {
			count++
		}
	}

	return count
}

func (h *GameHandler) Close() {
	close(h.Out)
}

func (h *GameHandler) Read() int {
	return <-h.In
}

func (h *GameHandler) Write(value int) {

	h.Out <- value
}

func (h *GameHandler) processOutput(x int, y int, tile int) {
	h.Screen[Point{x, y}] = tile

	switch tile {
	case Paddle:
		h.PaddleCoords = Point{x, y}
		break
	case Ball:
		h.BallCoords = Point{x, y}
		h.sendInput()
		break
	}
}

func (h *GameHandler) sendInput() {
	switch {
	case h.PaddleCoords.x < h.BallCoords.x:
		h.In <- 1
		break
	case h.PaddleCoords.x > h.BallCoords.x:
		h.In <- -1
	default:
		h.In <- 0
	}
}

func NewHandler() *GameHandler {
	return &GameHandler{
		In:     make(chan int, 1),
		Out:    make(chan int, 3),
		Screen: make(map[Point]int),
	}
}

func main() {
	input := goutils.ParseInput("input.txt")

	handler := NewHandler()

	input[0] = 2
	vm := intcode.NewVm(input, handler)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for x := range handler.Out {
			y := <-handler.Out
			tile := <-handler.Out

			handler.processOutput(x, y, tile)
		}
		wg.Done()
	}()

	vm.Run()

	handler.Close()
	wg.Wait()
	log.Println(handler.Screen[Point{-1, 0}])
}
