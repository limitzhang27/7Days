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

	r.GET("/hello", func(c *gee.Content) {
		c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Content) {
		c.Json(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	log.Fatalln(r.RUN(":9527"))
}
