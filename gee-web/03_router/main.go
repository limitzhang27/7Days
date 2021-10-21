package main

import (
	"gee"
	"log"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.Html(http.StatusOK, "<h1>Hello, World</h1>")
	})

	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you`re at %s\n", c.Param("name"), c.Path)
	})

	r.POST("/login", func(c *gee.Context) {
		c.Json(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.Json(http.StatusOK, gee.H{
			"filepath": c.Param("filpath"),
		})
	})

	log.Fatalln(r.RUN(":9527"))
}
