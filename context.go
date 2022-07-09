/*
===Context===
Context 随着每一个请求的出现而产生，请求的结束而销毁，和当前请求强相关的信息都应由 Context 承载。
因此，设计 Context 结构，扩展性和复杂性留在了内部，而对外简化了接口。
路由的处理函数，以及将要实现的中间件，参数都统一使用 Context 实例，
Context 就像一次会话的百宝箱，可以找到任何东西。

代码最开头，给map[string]interface{}起了一个别名gee.H，构建JSON数据时，显得更简洁。
Context目前只包含了http.ResponseWriter和*http.Request，另外提供了对 Method 和 Path 这两个常用属性的直接访问。
提供了访问Query和PostForm参数的方法。
提供了快速构造String/Data/JSON/HTML响应的方法。
*/

/*
===中间件===
中间件的定义与路由映射的 Handler 一致，处理的输入是Context对象。
插入点是框架接收到请求初始化Context对象后，允许用户使用自己定义的中间件做一些额外的处理，例如记录日志等，以及对Context进行二次加工。
另外通过调用(*Context).Next()函数，中间件可等待用户自己定义的 Handler处理结束后，再做一些额外的操作

之前的框架设计是这样的，当接收到请求后，匹配路由，该请求的所有信息都保存在Context中。
中间件也不例外，接收到请求后，应查找所有应作用于该路由的中间件，保存在Context中，依次进行调用。
*/

package giny

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	//
	Writer http.ResponseWriter
	Req    *http.Request

	//请求信息
	Path   string
	Method string
	params map[string]string

	//响应信息
	StatusCode int

	//中间件相关
	handlers []HandlerFunc
	index    int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1, //index如果初始化为0，去掉Next函数的首行自增，会卡死在第一个中间件，因为index一直为0
	}
}

//PostForm returns the first value for the named component of the query.
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//Query parses RawQuery and returns the corresponding values.
//It silently discards malformed value pairs.
//To check errors use ParseQuery.
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//Status sends an HTTP response header with the provided status code.
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//Set sets the header entries associated with key to the single element value.
//It replaces any existing values associated with key.
//The key is case insensitive;
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

//String sends an HTTP response with the provided code, format, values.
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//JSON sends an HTTP response with the provided code and data which was encoded in json.
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(&obj); err != nil {
		panic(err)
	}
}

//Data sends an HTTP response with the provided code and data.
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

//HTML sends an HTTP response with the provided code and HTML text.
func (c *Context) HTML(code int, HTMLText string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(HTMLText))
}

//Param 返回传入key的相应参数
func (c *Context) Param(key string) string {
	return c.params[key]
}

//用于中间件调用
func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		c.handlers[c.index](c)
		c.index++
	}
}
