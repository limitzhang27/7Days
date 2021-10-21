package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Content struct {
	// origin object
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
}

func newContent(w http.ResponseWriter, req *http.Request) *Content {
	return &Content{
		Path:   req.URL.Path,
		Method: req.Method,
		Req:    req,
		Writer: w,
		index:  -1,
	}
}

func (c *Content) Next() {
	c.index++
	s := len(c.handlers)
	// 这里使用遍历是一种兼容方式，有一些中间件不会手动执行Next，所以不能直接使用递归
	// 即使遇到有手动执行的中间件执行Next，对应的c.index也会发生改变，
	// 递归完回到上一层时，也会因为c.index的+1而跳过
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Content) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.Json(code, H{"message": err})
}

func (c *Content) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Content) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Content) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Content) Status(code int) {
	c.StatusCode = code
}

func (c *Content) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Content) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, _ = fmt.Fprintf(c.Writer, format, values...)
}

func (c *Content) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Content) Data(code int, data []byte) {
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

func (c *Content) Html(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, _ = c.Writer.Write([]byte(html))
}
