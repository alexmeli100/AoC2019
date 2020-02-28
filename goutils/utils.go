package goutils

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type IoChannel struct {
	In  chan int
	Out chan int
}

func (io *IoChannel) Read() (int, error) {
	return <-io.In, nil
}

func (io *IoChannel) Write(value int) {
	io.Out <- value
}

func (io *IoChannel) Close() {
	close(io.Out)
}

func NewIOCHan() *IoChannel {
	in := make(chan int, 2)
	out := make(chan int, 2)

	return &IoChannel{in, out}
}

// function for parsing intcode input
func ParseInput(path string) []int {
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

func Max(t []int) int {
	max := t[0]

	for _, v := range t {
		if v > max {
			max = v
		}
	}

	return max
}

// permutation implementation from https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
func Permutations(arr []int) [][]int {
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
