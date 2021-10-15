package gee

import (
	"log"
	"net/http"
)

type router struct {
	handles map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handles: map[string]HandlerFunc{}}
}

func getKey(method, pattern string) string {
	return method + "-" + pattern
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := getKey(method, pattern)
	r.handles[key] = handler
}

func (r *router) handle(c *Content) {
	log.Printf("HTTP method:%s pattern: %s", c.Method, c.Path)
	key := getKey(c.Method, c.Path)
	if handle, ok := r.handles[key]; ok {
		handle(c)
	} else {
		c.String(http.StatusNotFound, "404 Not Found path: %s", c.Path)
	}
}
