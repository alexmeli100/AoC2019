package main

import (
	"github.com/alexmeli100/AoC2019/day7/solution/intcode"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
)

type IoChannel struct {
	in  chan int
	out chan int
}

func (io *IoChannel) Read() int {
	return <-io.in
}

func (io *IoChannel) Write(value int) {
	io.out <- value
}

func (io *IoChannel) Close() {
	close(io.out)
}

func NewIOCHan() *IoChannel {
	in := make(chan int, 2)
	out := make(chan int, 2)

	return &IoChannel{in, out}
}

func main() {
	input := parseInput("input.txt")
	var wg sync.WaitGroup

	io := NewIOCHan()
	vm := intcode.NewVm(input, io)
	wg.Add(1)

	go func() {
		for output := range io.out {
			log.Println(output)
		}
		wg.Done()
	}()

	io.in <- 2
	vm.Run()
	close(io.out)
	wg.Wait()
}

func parseInput(path string) []int {
	bytes, err := ioutil.ReadFile(path)
	var res []int

	if err != nil {
		log.Fatal(err)
	}

	input := strings.Split(string(bytes), ",")

	for _, s := range input {
		i, err := strconv.Atoi(s)

		if err != nil {
			log.Fatal(err)
		}

		res = append(res, i)
	}

	return res
}