package msg

import "log"

type Producer struct {
	num   int
	topic string
	c     MessageChannel
}

func NewProducer(num int, topic string, chanSize int) *Producer {
	return &Producer{
		num:   num,
		topic: topic,
		c:     make(MessageChannel, chanSize),
	}
}

func (p *Producer) Start() {
	for i := 0; i < p.num; i++ {
		go func() {
			p.startProducer()
		}()
	}
}

func (p *Producer) startProducer() {
	go func() {
		for {
			select {
			case msg := <-p.c:
				data, err := EncodeMessage(msg)
				if err != nil {
					log.Println("[Producer]Encode Message err ", err)
					continue
				}
				log.Printf("[Producer] Send Data %v\n", msg)
				go func() {
					_ = p.udpSend(data)
				}()
			}
		}
	}()
}

func (p *Producer) udpSend(data []byte) error {
	log.Printf("[Producer]UPD Send data len(%d)", len(data))
	return nil
}

func (p *Producer) Send(message Message) {
	p.c <- message
}
