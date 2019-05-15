package api

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/pprof"
	"go.opencensus.io/trace"

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

var contextLoggerHTTP *log.Entry

const ginContextKeyMaster = "master"

func init() {
	contextLoggerHTTP = log.WithFields(log.Fields{
		"module": "http",
	})
}
func aa(b string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, c.MustGet(ginContextKeyMaster).(cue.MasterManager))
	}
}

func debug(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	master := c.MustGet(ginContextKeyMaster).(cue.MasterManager)
	debugData := struct {
		Config  *config.Server       `json:"config"`
		Hues    light.DiscoveredHues `json:"discovered_hues"`
		Version string               `json:"version"`
	}{
		Config:  config.GetServerConfig(ctx),
		Hues:    master.GetLightManager().GetDiscoveredHues(),
		Version: config.GetVersion(),
	}

	c.JSON(200, debugData)
}

func runCommands(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	var commands []string
	var responses []cue.Cue
	ctx, span := trace.StartSpan(ctx, "http: run commands")
	defer span.End()
	if err := c.ShouldBindJSON(&commands); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m := c.MustGet(ginContextKeyMaster).(cue.MasterManager)
	cs := m.GetDefaultCueStack()
	for _, eachCommand := range commands {
		ctx, span := trace.StartSpan(ctx, "http: run single Command")
		x, err := cue.CommandToCue(ctx, m, eachCommand)
		if err != nil {
			tracing.SetError(span, err)
			contextLoggerHTTP.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "command": eachCommand})
			span.End()
			return
		}
		x.Source.Input = cue.SourceInputAPI
		x.Source.Type = cue.SourceTypeCommand
		x.Source.Meta = eachCommand
		m.EnQueueCue(ctx, *x, cs)
		responses = append(responses, *x)
		span.End()
	}
	c.JSON(200, responses)
}
func getMaster(c *gin.Context) {
	c.JSON(200, c.MustGet(ginContextKeyMaster).(cue.MasterManager))
}

//createCue takes a JSON cue, and adds it to the default cuestack.
func createCue(c *gin.Context) {
	ctx := c.MustGet("ctx").(context.Context)
	ctx, span := trace.StartSpan(ctx, "createCue")
	defer span.End()
	var newCue cue.Cue
	if err := c.ShouldBindJSON(&newCue); err != nil {
		tracing.SetError(span, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m := c.MustGet(ginContextKeyMaster).(cue.MasterManager)
	stack := m.GetDefaultCueStack()
	newCue.Source.Input = cue.SourceInputAPI
	newCue.Source.Type = cue.SourceTypeJSON
	cue := m.EnQueueCue(ctx, newCue, stack)
	//span.SetTag("cue-id", cue.ID)
	c.JSON(200, cue)
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
		contextLoggerHTTP.Println("unable to encode image.")
	}

	c.Header("Content-Type", "image/jpeg")
	c.Writer.Write(buffer.Bytes())

}
func getLightInventory(c *gin.Context) {
	master := c.MustGet(ginContextKeyMaster).(cue.MasterManager)
	c.JSON(200, master.GetLightManager().GetLightsByName())
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
	wsTypeLightList   = "LIGHT_LIST"
	wsTypeCueList     = "CUE_MASTER"
	wsTypeLightStates = "LIGHT_STATES"
)

func wshandler(w http.ResponseWriter, r *http.Request, tickInterval time.Duration, master cue.MasterManager) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade ", err)
		return
	}

	go func() {
		for {
			//todo: only emit when things have changed
			conn.WriteJSON(&wsWrapper{Data: master.GetLightManager().GetLightsByName(), Type: wsTypeLightList})
			conn.WriteJSON(&wsWrapper{Data: master, Type: wsTypeCueList})
			conn.WriteJSON(&wsWrapper{Data: master.GetLightManager().GetAllStates(), Type: wsTypeLightStates})
			time.Sleep(tickInterval)
		}
	}()
}

func getRouter(ctx context.Context, master cue.MasterManager, testMode bool) *gin.Engine {
	if testMode {
		gin.SetMode(gin.TestMode)

	}
	httpConfig := config.GetServerConfig(ctx).Inputs.HTTP
	router := gin.New()
	router.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/_metrics"))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	//setup CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-JWT"}
	router.Use(cors.New(corsConfig))
	router.Use(tracing.GinMiddleware(ctx))

	//inject pointer to cuemaster into gin context
	router.Use(func(c *gin.Context) {
		c.Set(ginContextKeyMaster, master)
	})

	//register prometheus gin metrics middleware
	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	pprof.Register(router)

	//setup routes
	router.GET("/_metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/lights", getLightInventory)
	router.POST("cues", createCue)
	router.POST("commands", runCommands)
	router.GET("cuemaster", getMaster)
	router.GET("/ping", aa("ff"))
	router.GET("/hexfade/:from/:to", hexFade)
	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request, httpConfig.WSTickInterval, master)
	})
	router.GET("debug", debug)

	// router.Use(static.Serve("/", static.LocalFile("./ui/build", false)))
	return router
}

//ServeHTTP runs the gin server
func ServeHTTP(ctx context.Context, wg *sync.WaitGroup, master cue.MasterManager) {
	defer wg.Done()
	httpConfig := config.GetServerConfig(ctx).Inputs.HTTP
	if !httpConfig.Enabled {
		contextLoggerHTTP.Info("http is not enabled")
		return
	}

	//server
	srv := &http.Server{
		Addr:    httpConfig.Address,
		Handler: getRouter(ctx, master, false),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			contextLoggerHTTP.Fatalf("listen: %s\n", err)
		}
	}()
	<-ctx.Done()

	contextLoggerHTTP.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		contextLoggerHTTP.Fatal("Server Shutdown:", err)
	}
	contextLoggerHTTP.Println("Server exiting")
}
