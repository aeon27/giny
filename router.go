/*
将和路由相关的方法和结构提取出来,方便下一次对 router 的功能进行增强
*/

package giny

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       //不同请求方式对应的trie树根节点
	handlers map[string]HandlerFunc //不同pattern对应的handler
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	s := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range s {
		if item != "" {
			parts = append(parts, item)
			//有通配符*的路由，解析到此为止，只允许有一个 *
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	_, ok := r.roots[method]
	if !ok {
		//对未初始化的method对应的trie树进行初始化
		r.roots[method] = &node{}
	}

	//根据method对应的根节点和节点的insert方法构建路由
	parts := parsePattern(pattern)
	// log.Printf("parts: %v, len: %v\n", parts, len(parts))
	r.roots[method].insert(pattern, parts, 0)

	r.handlers[method+"-"+pattern] = handler
	// log.Printf("addRoute: %4s - %s\n", method, pattern)
}

//getRoute方法对传入的pattern进行参数解析，返回路由的叶结点和解析的参数map
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	params := make(map[string]string)
	searchParts := parsePattern(path)

	root, ok := r.roots[method]
	//方法对应的trie树不存在
	if !ok {
		return nil, nil
	}

	node := root.search(searchParts, 0)
	//将参数的解析结果存入params
	if node != nil {
		//参数解析——searchParts为请求的路径切片，parts为匹配到的路径的切片
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			//对于 : 动态路由的解析
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			//对于通配符 * 的解析
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
			}
		}

		return node, params
	}

	return nil, nil
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)

	if node != nil {
		c.params = params
		//第三个部分不能用c.Path，因为可能有参数，存在模糊匹配，所以应该用getRoute找到的node节点的pattern
		//r.handlers[c.Method+"-"+node.pattern](c)
		//将匹配到的handler添加到c.handlers的末尾，而不再直接执行
		c.handlers = append(c.handlers, r.handlers[c.Method+"-"+node.pattern])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}

	//依次执行c.handlers
	c.Next()
}
