package server

import (
	"net/http"
	"strings"
)

//返回一个Router实例
func NewRouter() *Router {
	return new(Router)
}

//路由结构体，包含一个记录方法、路径的map
type Router struct {
	RouteFunc    map[string]map[string]http.HandlerFunc
	RouteHandler map[string]http.Handler
}

//实现Handler接口，匹配方法以及路径
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h, ok := r.RouteFunc[req.Method][req.URL.String()]; ok {
		h(w, req)
	} else if h, ok := r.RouteHandler[req.URL.String()]; ok {
		http.Handle(req.URL.String(), h)
	}
}

//根据方法、路径将方法注册到路由
func (r *Router) HandleFunc(method, path string, f http.HandlerFunc) {
	method = strings.ToUpper(method)
	if r.RouteFunc == nil {
		r.RouteFunc = make(map[string]map[string]http.HandlerFunc)
	}
	if r.RouteFunc[method] == nil {
		r.RouteFunc[method] = make(map[string]http.HandlerFunc)
	}
	r.RouteFunc[method][path] = f
}

//根据方法、路径将方法注册到路由
func (r *Router) Handle(path string, f http.Handler) {
	if r.RouteHandler == nil {
		r.RouteHandler = make(map[string]http.Handler)
	}
	r.RouteHandler[path] = f
}
