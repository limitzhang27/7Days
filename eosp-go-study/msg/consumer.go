package msg

type Consumer struct {
	num   int
	topic string
	c     MessageChannel
}
