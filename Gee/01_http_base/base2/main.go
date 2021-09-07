package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct {
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		_, _ = fmt.Fprintf(w, "URL.PATH = [%q]\n", r.URL.Path)
		return
	case "/hello":
		data := ""
		for k, v := range r.Header {
			data += fmt.Sprintf("Header[%q] = [%q]\n", k, v)
		}
		_, _ = w.Write([]byte(data))
	default:
		_, _ = w.Write([]byte("404"))
	}
}

// 自定义路由器的web服务器
func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9527", engine))
}
