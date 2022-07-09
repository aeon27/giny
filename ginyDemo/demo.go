package main

import (
	"github.com/aeon27/giny"
	"net/http"
)

func main() {

	engine := giny.New()

	//为全局添加Logger中间件
	engine.Use(giny.Logger())

	engine.GET("/index", func(c *giny.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page<h1>")
	})

	v1 := engine.Group("/v1")
	{
		v1.GET("/", func(c *giny.Context) {
			c.HTML(http.StatusOK, "<h1>Hello v1!<h1>")
		})

		//用法：/hello?name=...
		v1.GET("/hello", func(c *giny.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := engine.Group("/v2")
	{
		//用法： /helo/...
		v2.GET("/hello/:name", func(c *giny.Context) {
			c.String(http.StatusOK, "heelo %s, you're at %s\n", c.PostForm("name"), c.Path)
		})
		v2.POST("/login", func(c *giny.Context) {
			c.JSON(http.StatusOK, giny.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	engine.Run(":9999")
}

// engine.GET("/", func(c *giny.Context) {
// 	// fmt.Fprintf(c.Writer, "URL.Path = %s\n", c.Path)
// 	c.HTML(http.StatusOK, "<h1>Hello giny</h1>")
// })

// engine.GET("/:name", func(c *giny.Context) {
// 	// for k, v := range c.Req.Header {
// 	// 	fmt.Fprintf(c.Writer, "Header[%q] = %q\n", k, v)
// 	// }
// 	c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
// })

// engine.GET("/hey", func(c *giny.Context) {
// 	// for k, v := range c.Req.Header {
// 	// 	fmt.Fprintf(c.Writer, "Header[%q] = %q\n", k, v)
// 	// }
// 	c.String(http.StatusOK, "hey %s, you are at %s\n", c.Query("name"), c.Path)
// })

// engine.GET("/hello/:name", func(c *giny.Context) {
// 	c.String(http.StatusOK, "hello %s, you are at %s\n", c.Param("name"), c.Path)
// })

// engine.GET("/assets/*filepath", func(c *giny.Context) {
// 	c.JSON(http.StatusOK, &giny.H{
// 		"filepath": c.Param("filepath"),
// 	})
// })

// engine.POST("/login", func(c *giny.Context) {
// 	c.JSON(http.StatusOK, &giny.H{
// 		"username": c.PostForm("username"),
// 		"password": c.PostForm("password"),
// 	})
// })
