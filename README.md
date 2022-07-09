# Giny

Giny是一个用Golang开发的，参考Gin实现的简易Web开发框架，支持基本的Web开发框架的功能，如静态路由、动态路由、路由分组控制、RESTful API和中间件扩展

# 使用帮助

ginyDemo文件夹中的demo.go文件提供了Giny的一些简单的使用范例，仅供参考 :D

# 快速开始

```
package main

import "github.com/aeon27/giny"

func main() {
	engine := giny.New()
	engine.Use(giny.Logger())//为全局使用logger中间件
	engine.GET("/", func(c *giny.Context) {
		c.HTML(200, "<h1>Hello Giny!~<h1>")
	})
	engine.Run(":9999") //打开浏览器输入 http://localhost:9999/ 可以看到运行效果
}
```
