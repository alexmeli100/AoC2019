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

func (io *IoChannel) Read() int {
	return <-io.In
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