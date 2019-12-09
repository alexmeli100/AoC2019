package main

import (
	"../solution/intcode"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
)

type Thruster struct {
	t    []int
	amps []*Amp
}

type Amp struct {
	io *IoChannel
	vm *intcode.IntCode
}

func NewAmp(tape []int) *Amp {
	io := NewIOCHan()
	vm := intcode.NewVm(tape, io)

	return &Amp{io, vm}
}

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

func NewThruster(t []int, tape []int) *Thruster {
	amps := make([]*Amp, len(t))

	for i := 0; i < len(t); i++ {
		amps[i] = NewAmp(tape)
	}

	return &Thruster{t, amps}
}

func (thrust *Thruster) connect() {
	last := len(thrust.t) - 1
	for i := 0; i < last; i++ {
		thrust.amps[i].io.in <- thrust.t[i]
		thrust.amps[i+1].io.in = thrust.amps[i].io.out
	}

	thrust.amps[last].io.in <- thrust.t[last]

	thrust.amps[0].io.in = thrust.amps[last].io.out
	thrust.amps[0].io.in <- thrust.t[0]
	thrust.amps[0].io.in <- 0
}

func (thrust *Thruster) signalLoop() int {
	thrust.connect()

	wg := sync.WaitGroup{}

	for i := range thrust.t {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			thrust.amps[i].vm.Run()
		}(i)
	}
	wg.Wait()
	return thrust.amps[len(thrust.amps)-1].vm.Last
}

func (thrust *Thruster) signal() int {
	out := 0

	for i, v := range thrust.t {
		thrust.amps[i].io.in <- v
		thrust.amps[i].io.in <- out
		thrust.amps[i].vm.Run()
		out = <-thrust.amps[i].io.out
	}

	return out
}

// change signal function to run part1 or part2
func signals(t []int, tape []int) []int {
	perms := permutations(t)
	out := make([]int, len(perms))

	for _, v := range perms {
		thrust := NewThruster(v, tape)
		res := thrust.signalLoop()
		out = append(out, res)
	}

	return out
}

func max(t []int) int {
	max := t[0]

	for _, v := range t {
		if v > max {
			max = v
		}
	}

	return max
}

// permutation implementation from https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
func permutations(arr []int) [][]int {
	var helper func([]int, int)
	var res [][]int

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
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

func main() {
	input := parseInput("input.txt")

	sigs := signals([]int{9, 8, 7, 6, 5}, input)
	out := max(sigs)

	log.Println(out)
}
