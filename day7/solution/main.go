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
	amps []*intcode.IntCode
}

func NewThruster(t []int, tape []int) *Thruster {
	amps := make([]*intcode.IntCode, len(t))

	for i := 0; i < len(t); i++ {
		amps[i] = intcode.NewVm(tape)
	}

	return &Thruster{t, amps}
}

func (thrust *Thruster) connect() {
	last := len(thrust.t) - 1
	for i := 0; i < last; i++ {
		thrust.amps[i].In <- thrust.t[i]
		thrust.amps[i+1].In = thrust.amps[i].Out
	}

	thrust.amps[last].In <- thrust.t[last]

	thrust.amps[0].In = thrust.amps[last].Out
	thrust.amps[0].In <- thrust.t[0]
	thrust.amps[0].In <- 0
}

func (thrust *Thruster) signalLoop() int {
	thrust.connect()

	wg := sync.WaitGroup{}

	for i := range thrust.t {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			thrust.amps[i].Run()
		}(i)
	}
	wg.Wait()
	return thrust.amps[len(thrust.amps)-1].Last
}

func (thrust *Thruster) signal() int {
	out := 0

	for i, v := range thrust.t {
		thrust.amps[i].In <- v
		thrust.amps[i].In <- out
		thrust.amps[i].Run()
		out = <-thrust.amps[i].Out
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
