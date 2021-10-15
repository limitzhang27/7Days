package gee

import (
	"fmt"
	"log"
)

type router struct {
	handles map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handles: map[string]HandlerFunc{}}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := fmt.Sprintf("%s-%s", method, pattern)
	r.handles[key] = handler
}

func (r *router) handle(c *Content) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handles[key]; ok {
		handler(c)
	}
}
