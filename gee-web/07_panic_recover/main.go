package main

import (
	"gee"
	"log"
	"net/http"
)

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello World\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"limigzhang"}
		c.String(http.StatusOK, names[100])
	})
	log.Fatal(r.Run(":9527"))
}
