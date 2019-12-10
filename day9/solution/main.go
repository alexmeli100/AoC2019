package main

import (
	"github.com/alexmeli100/AoC2019/day7/solution/intcode"
	"github.com/alexmeli100/AoC2019/goutils"
	"log"
	"sync"
)

func main() {
	input := goutils.ParseInput("input.txt")
	var wg sync.WaitGroup

	io := goutils.NewIOCHan()
	vm := intcode.NewVm(input, io)
	wg.Add(1)

	go func() {
		for output := range io.Out {
			log.Println(output)
		}
		wg.Done()
	}()

	io.In <- 2
	vm.Run()
	close(io.Out)
	wg.Wait()
}