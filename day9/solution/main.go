package main

import (
	"github.com/alexmeli100/AoC2019/day7/solution/intcode"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"
)

func main() {
	input := parseInput("input.txt")
	var wg sync.WaitGroup

	vm := intcode.NewVm(input)
	wg.Add(1)

	go func() {
		for output := range vm.Out {
			log.Println(output)
		}
		wg.Done()
	}()

	vm.In <- 1
	vm.Run()
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