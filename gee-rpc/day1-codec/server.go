package geerpc

import (
	"encoding/json"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int        // MagicNumber marks this`s a geerpc request
	CodecType   codec.Type // client amy chooses different Codec to encode body
}

var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

// Server represents an RPC Server.
type Server struct {
}

// NewServer returns a new Server.
func NewServer() *Server {
	return &Server{}
}

//DefaultServer is the default instance of *Server
var DefaultServer = NewServer()

// Accept accepts connections on the listener and servers requests
// for each incoming connection
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.ServerConn(conn)
	}
}

// ServerConn 处理TCP连接， 获取头部信息，并校验
func (server *Server) ServerConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: option error: ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x\n", opt.MagicNumber)
		return
	}

	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %x\n", opt.CodecType)
		return
	}
	server.serverCodec(f(conn))
}

// invalidRequest is a placeholder for response argv when error occurs
var invalidRequest = struct{}{}

// serverCodec 对conn读取数据，并解码
func (server *Server) serverCodec(cc codec.Codec) {
	sending := new(sync.Mutex) // make sure to send a complete response
	wg := new(sync.WaitGroup)  // wait until all request are handled
	// 一个连接中，会存在多个请求，使用for循环处理，直到发生错误
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break // it`s not possible to recover, so close the connection
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1) // 处理加一
		// 调用handle处理请求
		go server.handleRequest(cc, req, sending, wg)
	}
}

type request struct {
	h            *codec.Header // header of request
	argv, replyv reflect.Value // argv and replyv of request
}

//readRequestHeader 调用编码器，获取里面的头部信息，并返回
func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error: ", err)
		}
		return nil, err
	}
	return &h, nil
}

//readRequest 读取请求， 解析header,和body
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{
		h: h,
	}
	// TODO: now we dont`t know the type of request argv
	// day 1, just suppose it`s string
	req.argv = reflect.New(reflect.TypeOf(""))
	if err := cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err: ", err)
	}
	return req, nil
}

//handleRequest 处理请求，
func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// TODO, should call registered rpc methods to get the right replyv
	// day 1, just print argv and send a hello message
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

//sendResponse 发送请求
func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	// 单个连接中，有多个请求，需要多个返回，但是报文只能逐个发送，使用sync.Mutex锁保证
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error: ", err)
	}
}

func Accept(lis net.Listener) { DefaultServer.Accept(lis) }
