package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nickysemenza/hyperion/backend/cue"
	"github.com/nickysemenza/hyperion/backend/light"
)

func aa(b string) func(*gin.Context) {
	return func(c *gin.Context) {
		// c.JSON(200, gin.H{
		// 	"message": b,
		// })
		c.JSON(200, cue.CM)
	}
}
func getLightInventory(c *gin.Context) {
	c.JSON(200, light.Config)
}

//ServeHTTP runs the gin server
func ServeHTTP() {

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-JWT"}
	r.Use(cors.New(corsConfig))

	r.GET("/ping", aa("ff"))
	r.GET("/lights", getLightInventory)
	r.Run()
}
