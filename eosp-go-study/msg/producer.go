package main

import (
	"fmt"
	"log"
)

type Producer struct {
	topic string
	c     MessageChannel
	store *Store
}

func NewProducer(num int, topic string, chanSize int, store *Store) *Producer {
	return &Producer{
		topic: topic,
		c:     make(MessageChannel, chanSize),
		store: store,
	}
}

func (p *Producer) startProducer() {
	go func() {
		for {
			select {
			case msg := <-p.c:
				log.Printf("[Producer] Send Data %v\n", msg)
				go func() {
					_ = p.finalSend(msg)
				}()
			}
		}
	}()
}

func (p *Producer) finalSend(message Message) error {
	topic := message.Type
	data := message.Data
	byteData, err := EncodeMessage(data)
	if err != nil {
		return fmt.Errorf("[Producer]Encode Message err %s", err)
	}
	_ = p.store.Push(topic, byteData)
	return nil
}

func (p *Producer) Send(message Message) {
	p.c <- message
}
