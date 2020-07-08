package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"rdbviewer/shared"
)

const (
	HTTP_ADDR = ":7446"
	DB_ADDR   = "localhost:4002"
)

var BUILD_MODE = ""

type router struct {
	dbClient  shared.DatabaseClient
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags)

	h := &router{}

	// connect to database
	conn, err := grpc.Dial(DB_ADDR, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	h.dbClient = shared.NewDatabaseClient(conn)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	s := router.Group("/static/")
	if BUILD_MODE == "prod" {
		s.Use(func(ctx *gin.Context) {
			ctx.Header("Cache-Control", "public, max-age=1296000") // 15 days
		})
	}
	s.Static("/", "/build/static")

	onFrame := func(ctx *gin.Context) {
		ctx.File("/build/frame.html")
	}
	router.GET("/", onFrame)
	router.GET("/shows", onFrame)
	router.GET("/bootlegs", onFrame)
	router.GET("/shows/:id", onFrame)
	router.GET("/bootlegs/:id", onFrame)
	router.GET("/random", onFrame)
	router.GET("/stats", onFrame)
	router.GET("/dump", onFrame)

	router.POST("/data/front", h.onDataFront)
	router.POST("/data/search", h.onDataSearch)
	router.POST("/data/shows", h.onDataShows)
	router.POST("/data/show/:id", h.onDataShow)
	router.POST("/data/bootlegs", h.onDataBootlegs)
	router.POST("/data/bootleg/:id", h.onDataBootleg)
	router.POST("/data/random", h.onDataRandom)
	router.POST("/data/stats", h.onDataStats)
	router.GET("/dumpget", h.onPageDumpGet)
	router.HEAD("/dumpget", h.onPageDumpGet)

	log.Printf("[router] serving on %s", HTTP_ADDR)
	router.Run(HTTP_ADDR)
}
