package main

import (
	"bytes"
	"encoding/binary"
)

const TYPE_EMIAL = "EMAIL"
const TYPE_PHONE = "MSG"

type ByteView []byte

type Message struct {
	Type string
	Data interface{}
}

type MessageEmail struct {
	From    string
	To      string
	Content []byte
}

type MessagePhone struct {
	Phone   string
	Content []byte
}

type MessageChannel chan Message

func NewChannel(size int) MessageChannel {
	return make(chan Message, size)
}

func EncodeMessage(m interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, m); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
