package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Content) {
		// start time
		t := time.Now()
		// if a server error occured
		c.Fail(500, "Interface Service Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for grup v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger()) // global middleware
	r.GET("/", func(c *gee.Content) {
		c.Html(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gee.Content) {
			c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Param("name"))
		})
	}

	log.Println(r.RUN(":9527"))
}
