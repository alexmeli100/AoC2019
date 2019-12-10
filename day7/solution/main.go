package main

import (
	"../solution/intcode"
	"github.com/alexmeli100/AoC2019/goutils"
	"log"
	"sync"
)

type Thruster struct {
	t    []int
	amps []*Amp
}

type Amp struct {
	io *goutils.IoChannel
	vm *intcode.IntCode
}

func NewAmp(tape []int) *Amp {
	io := goutils.NewIOCHan()
	vm := intcode.NewVm(tape, io)

	return &Amp{io, vm}
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
		thrust.amps[i].io.In <- thrust.t[i]
		thrust.amps[i+1].io.In = thrust.amps[i].io.Out
	}

	thrust.amps[last].io.In <- thrust.t[last]

	thrust.amps[0].io.In = thrust.amps[last].io.Out
	thrust.amps[0].io.In <- thrust.t[0]
	thrust.amps[0].io.In <- 0
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
		thrust.amps[i].io.In <- v
		thrust.amps[i].io.In <- out
		thrust.amps[i].vm.Run()
		out = <-thrust.amps[i].io.Out
	}

	return out
}

// change signal function to run part1 or part2
func signals(t []int, tape []int) []int {
	perms := goutils.Permutations(t)
	out := make([]int, len(perms))

	for _, v := range perms {
		thrust := NewThruster(v, tape)
		res := thrust.signalLoop()
		out = append(out, res)
	}

	return out
}

func main() {
	input := goutils.ParseInput("input.txt")

	sigs := signals([]int{9, 8, 7, 6, 5}, input)
	out := goutils.Max(sigs)

	log.Println(out)
}
