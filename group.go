/*
分组控制(Group Control)是 Web 框架应提供的基础功能之一。
所谓分组，是指路由的分组。如果没有路由分组，我们需要针对每一个路由进行控制。
但是真实的业务场景中，往往某一组路由需要相似的处理。例如：
	以/post开头的路由匿名可访问。
	以/admin开头的路由需要鉴权。
	以/api开头的路由是 RESTful 接口，可以对接第三方平台，需要三方平台鉴权。

本实现也是以路由前缀来区分，并且支持分组的嵌套。
例如/post是一个分组，/post/a和/post/b可以是该分组下的子分组。
作用在/post分组上的中间件(middleware)，也都会作用在子分组，子分组还可以应用自己特有的中间件。
*/

package giny

import "log"

type RouterGroup struct {
	prefix      string
	router      *router
	middlewares []HandlerFunc
}

//根路由分组
func newRootGroup() *RouterGroup {
	return &RouterGroup{
		prefix: "",
		router: newRouter(),
	}
}

//路由分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		router: group.router,
	}

	return newGroup
}

func (group *RouterGroup) addRoute(method, subPattern string, handler HandlerFunc) {
	pattern := group.prefix + subPattern
	group.router.addRoute(method, pattern, handler)
	log.Printf("Group route: %4s - %4s\n", method, pattern)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//暴露给用户的接口，供用户为某个分组添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
