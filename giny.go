package giny

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	groups []*RouterGroup
}

//创建一个Engine实例
func New() *Engine {
	group := newRootGroup()
	return &Engine{
		RouterGroup: group,
		groups:      []*RouterGroup{group},
	}
}

//engine的Group调用RouterGroup的Group方法，将返回的newGroup加入groups列表
func (engine *Engine) Group(prefix string) *RouterGroup {
	newGroup := engine.RouterGroup.Group(prefix)
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		//如果一个请求的前缀包含了某个分组，那么这个分组的中间件就会被应用到这个请求上
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	//Context 随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应由 Context 承载。
	c := newContext(w, req)
	//因为每个请求对应一个context，所以将要应用的中间件加到context.middlewares里即可
	c.handlers = middlewares
	engine.router.handle(c)
}

//运行Engine实例
func (engine *Engine) Run(addr string) error {
	log.Printf("Giny is Listening and serving on %s\n", addr)
	return http.ListenAndServe(addr, engine)
}
