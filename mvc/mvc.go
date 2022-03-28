package mvc

import (
	"net/http"
)

type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	middlewares []HandlerFunc
	router      *router
}

// New is the constructor of gee.Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	c.handlers = engine.middlewares
	engine.router.handle(c)
}

func (engine *Engine) Use(middlewares ...HandlerFunc) {
	engine.middlewares = append(engine.middlewares, middlewares...)
}
