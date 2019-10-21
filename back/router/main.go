package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
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

type Router struct {
	templates map[string]*template.Template
	dbClient  shared.DatabaseClient
}

func main() {
	pprofMux := http.DefaultServeMux
	go func() {
		(&http.Server{
			Addr:    ":9999",
			Handler: pprofMux,
		}).ListenAndServe()
	}()
	http.DefaultServeMux = http.NewServeMux()

	rand.Seed(time.Now().UnixNano())

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)

	h := &Router{}

	// connect to database
	conn, err := grpc.Dial(DB_ADDR, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	h.dbClient = shared.NewDatabaseClient(conn)

	h.templates = TplLoadAll("/build/template")

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// populate router
	s := r.Group("/static/")
	if BUILD_MODE == "prod" {
		s.Use(func(c *gin.Context) {
			c.Header("Cache-Control", "public, max-age=1296000") // 15 days
		})
	}
	s.Static("/", "/build/static")
	r.GET("/", h.onPageFront)
	r.POST("/data/search", h.onDataSearch)
	r.GET("/shows", h.onPageShows)
	r.POST("/data/shows", h.onDataShows)
	r.GET("/bootlegs", h.onPageBootlegs)
	r.POST("/data/bootlegs", h.onDataBootlegs)
	r.GET("/shows/:id", h.onPageShow)
	r.GET("/bootlegs/:id", h.onPageBootleg)
	r.GET("/random", h.onPageRandom)
	r.GET("/stats", h.onPageStats)
	r.GET("/dump", h.onPageDump)
	r.GET("/dumpget", h.onPageDumpGet)
	r.HEAD("/dumpget", h.onPageDumpGet)

	log.Printf("serving router on %s", HTTP_ADDR)
	r.Run(HTTP_ADDR)
}
