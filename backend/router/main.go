package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"rdbviewer/shared"
)

const (
	HTTP_ADDR = ":7446"
	DB_ADDR   = "localhost:4002"
)

var BUILD_MODE string

type router struct {
	templates templateMap
	dbClient  shared.DatabaseClient
}

func main() {
	pprofInit()

	rand.Seed(time.Now().UnixNano())

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)

	h := &router{}

	// connect to database
	conn, err := grpc.Dial(DB_ADDR, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	h.dbClient = shared.NewDatabaseClient(conn)

	h.templates = templateLoadAll("/build/templates")

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	s := router.Group("/static/")
	if BUILD_MODE == "prod" {
		s.Use(func(c *gin.Context) {
			c.Header("Cache-Control", "public, max-age=1296000") // 15 days
		})
	}
	s.Static("/", "/build/static")

	router.GET("/", h.onPageFront)
	router.POST("/data/search", h.onDataSearch)
	router.GET("/shows", h.onPageShows)
	router.POST("/data/shows", h.onDataShows)
	router.GET("/bootlegs", h.onPageBootlegs)
	router.POST("/data/bootlegs", h.onDataBootlegs)
	router.GET("/shows/:id", h.onPageShow)
	router.GET("/bootlegs/:id", h.onPageBootleg)
	router.GET("/random", h.onPageRandom)
	router.GET("/stats", h.onPageStats)
	router.GET("/dump", h.onPageDump)
	router.GET("/dumpget", h.onPageDumpGet)
	router.HEAD("/dumpget", h.onPageDumpGet)

	log.Printf("[router] serving on %s", HTTP_ADDR)
	router.Run(HTTP_ADDR)
}
