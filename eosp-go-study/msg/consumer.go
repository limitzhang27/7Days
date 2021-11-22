package main

type Consumer struct {
	topic string
	c     MessageChannel
	store *Store
}

func NewConsumer(num int, topic string, chanSize int, store *Store) *Consumer {
	return &Consumer{
		topic: topic,
		c:     make(MessageChannel, chanSize),
		store: store,
	}
}
