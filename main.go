package main

import (
	"github.com/gin-gonic/gin"
	"socialMedia/router"
)

func main() {

	// Create a new Gin router
	r := gin.Default()

	router.Page(r)

	// HTTP handler for creating a new post
	r.POST("/posts", func(c *gin.Context) {
		// Decode JSON request body into a Post struct

	})

	// Start the HTTP server
	r.Run(":8180")

}
