package main

import (
	"fmt"
	"log"
	"net/http"
)

// 一个简单的web服务器
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	log.Fatalln(http.ListenAndServe(":9527", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	for s, strings := range r.Header {
		_, _ = fmt.Fprintf(w, "Header[%q] = [%q]\n", s, strings)
	}
}
