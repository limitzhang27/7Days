package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

func (g *GobCodec) Close() error {
	return g.conn.Close()
}

func (g *GobCodec) ReadHeader(header *Header) error {
	return g.dec.Decode(header)
}

func (g *GobCodec) ReadBody(i interface{}) error {
	return g.dec.Decode(i)
}

func (g *GobCodec) Write(header *Header, i interface{}) (err error) {
	defer func() {
		_ = g.buf.Flush() // 缓冲池刷新，会相应的写到conn中
		if err != nil {
			_ = g.Close()
		}
	}()

	// 使用enc加密header, 因为enc绑定的是缓冲池buf, 加密后的数据会保存在缓冲池
	if err := g.enc.Encode(header); err != nil {
		log.Println("rpc codec: gob error encoding header: ", err)
		return err
	}
	if err := g.enc.Encode(i); err != nil {
		log.Println("rpc codec: gob error encoding body: ", err)
		return err
	}
	return nil
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn) // 给IO套接字加一个缓冲池,相当于IO套接字就是缓冲池输出的地方
	// buf的输入有encoder编码之后保存在缓冲池中，输入在IO
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}
