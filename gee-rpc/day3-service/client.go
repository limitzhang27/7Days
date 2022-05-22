package geerpc

import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"sync"
)

type Call struct {
	Seq           uint64
	ServiceMethod string      // format "<service>.<method>"
	Args          interface{} // arguments to the function
	Reply         interface{} // reply from the function
	Error         error       // function error occurs, it will be set
	Done          chan *Call  // Strobes when call is complete
}

func (call *Call) done() {
	call.Done <- call
}

// Client represents an RPC client
// There may be multiple outstanding Calls associated
// multiple goroutines simultaneously.
type Client struct {
	cc       codec.Codec // 消息编码器，服务器类似，用来序列化要发哦少年宫出去的请求，以及反序列化接收到的相应
	opt      *Option
	sending  sync.Mutex   // 保证请求的有序发送，即防止出现多个请求报文混淆
	header   codec.Header // 每个请求的消息头，header 只有在发送请求时才需要，而请求发哦是那个是互斥的，因此每个客户端都需要
	mu       sync.Mutex   // protect following
	seq      uint64
	pending  map[uint64]*Call // 存储未处理完的请求，键是编号， 值是 Call 实例
	closing  bool             // user has called Close
	shutdown bool             // server has told us stop
}

var ErrShutdown = errors.New("connection is shut down")

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing {
		return ErrShutdown
	}
	c.closing = true
	return c.cc.Close()
}

func (c *Client) IsAvailable() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return !c.shutdown && !c.closing
}

//registerCall 注册调用，将调用器注册到client中
func (c *Client) registerCall(call *Call) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing || c.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = c.seq
	c.pending[call.Seq] = call
	c.seq++
	return call.Seq, nil
}

//removeCall 从pending中取出call
func (c *Client) removeCall(seq uint64) *Call {
	c.mu.Lock()
	defer c.mu.Unlock()
	call := c.pending[seq]
	delete(c.pending, seq)
	return call
}

//terminateCall 服务端或者客户端发生错误时调用， 将 shutdown 设置为 true， 且将错误信息通知所有 pending 状态的 call
func (c *Client) terminateCall(err error) {
	c.sending.Lock()
	defer c.sending.Unlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.shutdown = true
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
}

//receive 客户端接受消息，接收到之后需要将调用器标记已完成
func (c *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = c.cc.ReadHeader(&h); err != nil {
			break
		}
		call := c.removeCall(h.Seq)
		switch {
		case call == nil:
			// it usually means that Write partially failed
			// and call was already removed.
			err = c.cc.ReadBody(nil)
		case h.Error != "":
			call.Error = fmt.Errorf(h.Error)
			err = c.cc.ReadHeader(nil)
			call.done()
		default:
			err = c.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	// error occurs, so terminateCalls pending calls
	c.terminateCall(err)
}

//send 发送消息
func (c *Client) send(call *Call) {
	// make sure that the client will send a complete request
	c.sending.Lock()
	defer c.sending.Unlock()

	// 调用器注册到客户端中
	seq, err := c.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	// 将调用器头部信息拷贝到client中
	c.header.ServiceMethod = call.ServiceMethod
	c.header.Seq = seq
	c.header.Error = ""

	// encode and send the request
	if err := c.cc.Write(&c.header, call.Args); err != nil {
		call := c.removeCall(seq)
		// call may be nil, it usually means that Write partially failed.
		// client has received the response and handled
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

//Go 真正的发起一次调用
func (c *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	c.send(call)
	return call
}

func (c *Client) Call(serviceMethod string, args, reply interface{}) error {
	call := <-c.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

var _ io.Closer = (*Client)(nil)

//NewClient 新建一个客户端， 先见Option以json形式发送给服务端
func NewClient(conn net.Conn, opt *Option) (*Client, error) {
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		log.Println("rpc client: codec error: ", err)
		return nil, err
	}

	// send option with server
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error: ", err)
		_ = conn.Close()
		return nil, err
	}
	return newClientCodec(f(conn), opt), nil
}

//newClientCodec 真正创建客户端，携带编码器
func newClientCodec(cc codec.Codec, opt *Option) *Client {
	client := &Client{
		seq:     1, // seq starts with 1, 0 means invalid call
		cc:      cc,
		opt:     opt,
		sending: sync.Mutex{},
		pending: make(map[uint64]*Call),
	}
	// 开启协程，接受服务端传过来的消息
	go client.receive()
	return client
}

//parseOptions
func parseOptions(opts ...*Option) (*Option, error) {
	// if opts is nil os pass nil as parameter
	if len(opts) == 0 || opts[0] == nil {
		return nil, errors.New("number of options is more than 1")
	}
	opt := opts[0]
	opt.MagicNumber = DefaultOption.MagicNumber
	if opt.CodecType == "" {
		opt.CodecType = DefaultOption.CodecType
	}
	return opt, nil
}

//Dial 与服务器建立连接，并新建一个客户端
func Dial(network, address string, opts ...*Option) (client *Client, err error) {
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	// close the connection if client is nil
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()

	return NewClient(conn, opt)
}
