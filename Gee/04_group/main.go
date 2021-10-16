package main

import (
	"gee"
	"log"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Content) {
		c.Html(http.StatusOK, "<h1>Hello, World</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Content) {
			c.Html(http.StatusOK, "<h1>Hello World</h1>")
		})

		v1.GET("/hello", func(c *gee.Content) {
			c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *gee.Content) {
			c.String(http.StatusOK, "hello %s, your`re at %s\n", c.Param("name"), c.Path)
		})

		v2.POST("login", func(c *gee.Content) {
			c.Json(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	log.Fatalln(r.RUN(":9527"))
}
