package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func getKey(method, pattern string) string {
	return method + "-" + pattern
}

// only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由，规则pattern对应handler, 再往roots中插入结点
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := getKey(method, pattern)
	parts := parsePattern(pattern)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 获取路由的结点node，并解析里面的参数
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 拆分路径转
	params := make(map[string]string) // 路径中的参数
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0) // 查找对应结点
	if n != nil {
		parts := parsePattern(n.pattern) // 解析路由规则
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index] // 从规则中找到可变参数，获取对应路径中的值
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
	}
	return n, params
}

func (r *router) handle(c *Content) {
	log.Printf("HTTP method:%s pattern: %s", c.Method, c.Path)
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := getKey(c.Method, n.pattern)
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 Not Found path: %s", c.Path)
	}
}
