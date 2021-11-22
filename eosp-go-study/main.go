package main

import (
	"fmt"
	"time"
)

type Message struct {
	Msg string
}

var c = make(chan Message, 100)

func main() {

	for i := 0; i < 5; i++ {
		go func(index int) {
			for j := 0; j < 10; j++ {
				producer(Message{Msg: fmt.Sprintf("[%d] send msg: %d", index, j)})
				time.Sleep(time.Second)
			}
		}(i)
	}

	for i := 0; i < 5; i++ {
		go consumer()
	}

	time.Sleep(100 * time.Second)
}

func producer(message Message) {
	c <- message
}

func consumer() {
	for {
		select {
		case msg := <-c:
			fmt.Printf("[consumer] get msg : %s\n", msg.Msg)
			time.Sleep(1 * time.Second)
		}
	}
}
