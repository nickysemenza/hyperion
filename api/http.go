package api

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"net/http"
	"os"
	"os/signal"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"

	"github.com/nickysemenza/hyperion/util/tracing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	colorful "github.com/lucasb-eyer/go-colorful"
	ginprometheus "github.com/mcuadros/go-gin-prometheus"
	"github.com/nickysemenza/hyperion/core/config"
	"github.com/nickysemenza/hyperion/core/cue"
	"github.com/nickysemenza/hyperion/core/light"
	"github.com/nickysemenza/hyperion/util/color"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func aa(b string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, cue.GetCueMaster())
	}
}

func debug(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	debugData := struct {
		Config  *config.Server       `json:"config"`
		Hues    light.DiscoveredHues `json:"discovered_hues"`
		Version string               `json:"version"`
	}{
		Config:  config.GetServerConfig(ctx),
		Hues:    light.GetDiscoveredHues(ctx),
		Version: config.GetVersion(),
	}

	c.JSON(200, debugData)
}

func runCommands(c *gin.Context) {
	var commands []string
	var responses []cue.Cue
	if err := c.ShouldBindJSON(&commands); err == nil {
		cs := cue.GetCueMaster().GetDefaultCueStack()
		for _, eachCommand := range commands {

			if x, err := cue.NewFromCommand(eachCommand); err != nil {
				log.Println(err)
			} else {
				cs.EnQueueCue(*x)
				responses = append(responses, *x)
			}
		}

		c.JSON(200, responses)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

}
func getCueMaster(c *gin.Context) {
	c.JSON(200, cue.GetCueMaster())
}

//createCue takes a JSON cue, and adds it to the default cuestack.
func createCue(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	span, ctx := opentracing.StartSpanFromContext(ctx, "createCue")
	defer span.Finish()
	var newCue cue.Cue
	if err := c.ShouldBindJSON(&newCue); err == nil {
		stack := cue.GetCueMaster().GetDefaultCueStack()
		cue := stack.EnQueueCue(newCue)
		span.SetTag("cue-id", cue.ID)
		c.JSON(200, cue)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

//hexFade returns an image representing the fade from one hex val to another
// NOTE: hex values must be given without the pound
func hexFade(c *gin.Context) {

	blocks := 20
	blockw := 40
	img := image.NewRGBA(image.Rect(0, 0, blocks*blockw, blockw))

	c1, _ := colorful.Hex("#" + c.Param("from"))
	c2, _ := colorful.Hex("#" + c.Param("to"))

	rgb1 := color.GetRGBFromColorful(c1)
	rgb2 := color.GetRGBFromColorful(c2)

	for i := 0; i < blocks; i++ {
		eachColor := rgb1.GetInterpolatedFade(rgb2, i, blocks)
		draw.Draw(img, image.Rect(i*blockw, 0, (i+1)*blockw, blockw), &image.Uniform{eachColor.AsColorful()}, image.ZP, draw.Src)
	}

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	c.Header("Content-Type", "image/jpeg")
	c.Writer.Write(buffer.Bytes())

}
func getLightInventory(c *gin.Context) {
	c.JSON(200, light.GetLights())
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsWrapper struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

const (
	wsTypeLightList = "LIGHT_LIST"
	wsTypeCueList   = "CUE_MASTER"
)

func wshandler(w http.ResponseWriter, r *http.Request, tickInterval time.Duration) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade ", err)
		return
	}

	go func() {
		for {

			conn.WriteJSON(&wsWrapper{Data: light.GetLights(), Type: wsTypeLightList})
			conn.WriteJSON(&wsWrapper{Data: cue.GetCueMaster(), Type: wsTypeCueList})
			time.Sleep(tickInterval)
		}
	}()

	// for {
	// 	t, msg, err := conn.ReadMessage()
	// 	if err != nil {
	// 		break
	// 	}
	// 	conn.WriteMessage(t, msg)
	// }
}

//ServeHTTP runs the gin server
func ServeHTTP(ctx context.Context) {
	httpConfig := config.GetServerConfig(ctx).Inputs.HTTP
	if !httpConfig.Enabled {
		log.Info("http is not enabled")
		return
	}
	router := gin.Default()

	//setup CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-JWT"}
	router.Use(cors.New(corsConfig))
	router.Use(tracing.GinMiddleware(ctx))

	//register prometheus gin metrics middleware
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	//setup routes
	router.GET("/_metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/lights", getLightInventory)
	router.POST("cues", createCue)
	router.POST("commands", runCommands)
	router.GET("cuemaster", getCueMaster)
	router.GET("/ping", aa("ff"))
	router.GET("/hexfade/:from/:to", hexFade)
	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request, httpConfig.WSTickInterval)
	})
	router.GET("debug", debug)

	//server
	srv := &http.Server{
		Addr:    httpConfig.Address,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//block until graceful shutdown signal
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
