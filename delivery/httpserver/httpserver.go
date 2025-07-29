package httpserver

import "github.com/gin-gonic/gin"

func Server() {
	r := gin.New() // Creates a router without default middleware

	// Use Logger and Recovery middleware globally
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, Gin!"})
	})
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	err := r.Run(":8080")
	if err != nil {
		return
	} // Start the server

}
