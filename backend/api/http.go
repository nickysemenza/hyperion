package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nickysemenza/hyperion/backend/cue"
	"github.com/nickysemenza/hyperion/backend/light"
)

func aa(b string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, cue.CM)
	}
}
func getLightInventory(c *gin.Context) {
	c.JSON(200, light.Config)
}

//ServeHTTP runs the gin server
func ServeHTTP() {

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-JWT"}
	router.Use(cors.New(corsConfig))

	router.GET("/ping", aa("ff"))
	router.GET("/lights", getLightInventory)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
