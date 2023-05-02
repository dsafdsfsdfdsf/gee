package gee

import (
	"log"
	"net/http"
)

type router struct {
	// 保管路由的分发
	handlers map[string]HandlerFunc
}

// newRouter 对router 内的mapping 进行初始化分配内存
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// addRoute 对mapping结构进行路由的添加
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		// handler是一个函数类型
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
