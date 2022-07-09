/*
中间件是什么？以及为什么要有中间件
中间件(middlewares)，简单说，就是非业务的技术类组件。
Web 框架本身不可能去理解所有的业务，因而不可能实现所有的功能。
因此，框架需要有一个插口，允许用户自己定义功能，嵌入到框架中，仿佛这个功能是框架原生支持的一样。
*/

package giny

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		// c.String(http.StatusOK, "do something...")
		c.Next()
		log.Printf("|| handle: %4s - %25s || Duration: %10s ||\n", c.Method, c.Path, time.Since(t))
	}
}
