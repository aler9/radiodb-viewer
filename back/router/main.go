package main

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"rdbviewer/shared"
	"time"

	_ "net/http/pprof"
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

	exe, _ := os.Executable()
	os.Chdir(filepath.Dir(exe))

	rand.Seed(time.Now().UnixNano())

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)

	h := &Router{
		templates: make(map[string]*template.Template),
	}

	// load templates
	func() {
		tplRoot := "./template"
		files, err := ioutil.ReadDir(tplRoot)
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			name := f.Name()
			name = name[:len(name)-len(filepath.Ext(name))]
			h.templates[name] = TplLoad(filepath.Join(tplRoot, f.Name()))
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// populate router
	func() {
		s := r.Group("/static/")
		if BUILD_MODE == "prod" {
			s.Use(func(c *gin.Context) {
				c.Header("Cache-Control", "public, max-age=1296000") // 15 days
			})
		}
		s.Static("/", "./static")

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
	}()

	// connect to database
	conn, err := grpc.Dial(DB_ADDR, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	h.dbClient = shared.NewDatabaseClient(conn)

	h.Log("serving router on %s", HTTP_ADDR)
	r.Run(HTTP_ADDR)
}

func (*Router) Log(text string, args ...interface{}) {
	log.Printf(text, args...)
}
