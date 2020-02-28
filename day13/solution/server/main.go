package main

import (
	"fmt"
	"github.com/alexmeli100/AoC2019/day7/solution/intcode"
	"github.com/alexmeli100/AoC2019/goutils"
	"github.com/bitly/go-nsq"
	"log"
	"strconv"
	"sync"
)

type MessageHandler struct {
	In  chan int
	Out chan int
}

func (h *MessageHandler) Close() {
	close(h.Out)
}

func (h *MessageHandler) Read() (int, error) {
	input := <-h.In

	return input, nil
}

func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	input, err := strconv.Atoi(string(m.Body))

	if err != nil {
		log.Fatal(err)
	}

	log.Println("received", input)
	h.In <- input

	return nil
}

func (h *MessageHandler) Write(value int) {
	h.Out <- value
}

func NewHandler() *MessageHandler {
	return &MessageHandler{
		In:  make(chan int, 1),
		Out: make(chan int, 3),
	}
}

func main() {
	input := goutils.ParseInput("input.txt")
	cfg := nsq.NewConfig()

	consumer, err := nsq.NewConsumer("game_input", "game", cfg)
	if err != nil {
		log.Fatal(err)
	}

	config := nsq.NewConfig()
	p, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Fatal(err)
	}

	handler := NewHandler()

	consumer.AddHandler(handler)
	err = consumer.ConnectToNSQLookupd("localhost:4161")

	if err != nil {
		log.Fatal(err)
	}

	vm := intcode.NewVm(input, handler)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for x := range handler.Out {

			y := <-handler.Out
			tile := <-handler.Out
			output := []byte(fmt.Sprintf("%d %d %d", x, y, tile))

			err = p.Publish("game_output", output)
			if err != nil {
				log.Fatal(err)
			}
		}
		wg.Done()
	}()

	vm.Run()
	handler.Close()
	wg.Wait()
	consumer.Stop()
	p.Stop()
}
