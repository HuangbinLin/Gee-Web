package main

/*
(1) global middleware Logger
$ curl http://localhost:9999/
<h1>Hello Gee</h1>

>>> log
2019/08/17 01:37:38 [200] / in 3.14µs
*/

/*
(2) global + group middleware
$ curl http://localhost:9999/v2/hello/geektutu
{"message":"Internal Server Error"}

>>> log
2019/08/17 01:38:48 [200] /v2/hello/geektutu in 61.467µs for group v2
2019/08/17 01:38:48 [200] /v2/hello/geektutu in 281µs
*/

// 分组的作用，相同的前缀有相同的功能，形同的前缀可能需要相同的前处理(中间件)
// 相同的前处理后，需要各自的handler的处理，处理后可能需要相同的后处理(中间件)
// 使用USE把中间件绑定在Group上
// 在handler中，检查当前请求与所有Group的前缀关系，然后把对应的中间件绑定到context中
// 默认的全局Group开头为/，所有会自动绑定他的handler
// 处理后，会把

import (
	"log"
	"net/http"
	"time"

	"gee"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger()) // global midlleware
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
