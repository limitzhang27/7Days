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
	// response
	StatusCode int
}

func newContent(w http.ResponseWriter, req *http.Request) *Content {
	return &Content{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
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

func (c *Content) Data(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, _ = c.Writer.Write([]byte(html))
}
