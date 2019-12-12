package main

import (
	"bufio"
	"fmt"
	"github.com/alexmeli100/AoC2019/day7/solution/intcode"
	"github.com/alexmeli100/AoC2019/goutils"
	"github.com/alexmeli100/go-netcat/server"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type IoTCP struct {
	conn *bufio.Reader
	Out  chan int
}

func (io *IoTCP) Read() int {
	log.Println("Expecting input...")
	input, err := io.conn.ReadString('\n')

	if err != nil {
		log.Fatal("Error reading from connection")
	}

	log.Println("Received input: ", input)
	res, err := strconv.Atoi(strings.TrimSpace(input))

	if err != nil {
		log.Fatal("Error parsing input")
	}

	return res
}

func (io *IoTCP) Write(value int) {
	io.Out <- value
}

func (io *IoTCP) Close() {
	close(io.Out)
}

func NewIo(conn net.Conn) *IoTCP {
	return &IoTCP{
		conn: bufio.NewReader(conn),
		Out:  make(chan int, 2),
	}
}

type IntHandler struct {
	tape []int
}

func (i IntHandler) Handle(conn net.Conn) {
	io := NewIo(conn)
	vm := intcode.NewVm(i.tape, io)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		for color := range io.Out {
			dir := <-io.Out
			output := []byte(fmt.Sprintf("%d %d\n", color, dir))
			log.Println("Sending ", output)
			conn.Write(output)
		}
		wg.Done()
	}()

	vm.Run()
	io.Close()
	wg.Wait()
	conn.Write([]byte("DONE\n"))
}

func main() {
	input := goutils.ParseInput("input.txt")
	s, err := server.NewServer("localhost:8080")

	if err != nil {
		log.Fatal("Error creating server")
	}

	s.Handle(IntHandler{input})

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	s.Stop()
}
