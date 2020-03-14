package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"rdbviewer/defs"
	"rdbviewer/shared"
)

const (
	DB_UPLOADER_ADDR = ":4001"
	DB_API_ADDR      = ":4002"
	DB_FILE          = "/data/radiodb.json"
)

type database struct {
	mutex sync.RWMutex

	data            defs.RadioOut
	showKeywords    map[string]map[string]struct{}
	bootlegKeywords map[string]map[string]struct{}
	perYearShows    map[uint32]int32
	perYearBootlegs map[uint32]int32
	perYearSize     map[uint32]uint64
}

func main() {
	pprofInit()

	rand.Seed(time.Now().UnixNano())

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)

	db := &database{}

	db.load()

	go db.serveUploader()
	db.serveApi()
}

func (db *database) load() {
	if _, err := os.Stat(DB_FILE); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}

		log.Printf("[db] %s does not exist and is not loaded")
		return
	}

	db.data = defs.RadioOut{}
	mustImportJson(DB_FILE, &db.data)

	db.showKeywords = make(map[string]map[string]struct{})
	for _, s := range db.data.Shows {
		db.showKeywords[s.Id] = make(map[string]struct{})
		pushText := func(text string) {
			for word := range getTextKeywords(text, 2) {
				db.showKeywords[s.Id][word] = struct{}{}
			}
		}
		pushText(shared.LabelArtist(s))
		pushText(s.City)
		pushText(s.CountryCode)
		pushText(shared.LabelCountry(s))
		d, _ := time.Parse("2006-01-02", s.Date)
		pushText(d.Format("01 02 2006 January"))
	}

	db.bootlegKeywords = make(map[string]map[string]struct{})
	for _, b := range db.data.Bootlegs {
		db.bootlegKeywords[b.Id] = make(map[string]struct{})
		pushText := func(text string) {
			for word := range getTextKeywords(text, 2) {
				db.bootlegKeywords[b.Id][word] = struct{}{}
			}
		}

		pushText(b.Name)

		s := db.data.Shows[b.Show]
		pushText(shared.LabelArtist(s))
		pushText(s.City)
		pushText(s.CountryCode)
		pushText(shared.LabelCountry(s))
		d, _ := time.Parse("2006-01-02", s.Date)
		pushText(d.Format("01 02 2006 January"))
	}

	db.perYearShows = make(map[uint32]int32)
	db.perYearBootlegs = make(map[uint32]int32)
	db.perYearSize = make(map[uint32]uint64)
	for i := uint32(1992); i <= shared.Atoui32(time.Now().Format("2006")); i++ {
		db.perYearShows[i] = 0
		db.perYearBootlegs[i] = 0
		db.perYearSize[i] = 0
	}
	for _, s := range db.data.Shows {
		d, _ := time.Parse("2006-01-02", s.Date)
		year := shared.Atoui32(d.Format("2006"))
		db.perYearShows[year]++
	}
	for _, b := range db.data.Bootlegs {
		s := db.data.Shows[b.Show]
		d, _ := time.Parse("2006-01-02", s.Date)
		year := shared.Atoui32(d.Format("2006"))
		db.perYearBootlegs[year]++
		db.perYearSize[year] += b.Size
	}
}

func (db *database) serveUploader() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.POST("/upload", db.onUpload)

	log.Printf("[db] serving uploader on %s", DB_UPLOADER_ADDR)
	router.Run(DB_UPLOADER_ADDR)
}

func (db *database) serveApi() {
	server := grpc.NewServer()
	shared.RegisterDatabaseServer(server, db)

	listener, err := net.Listen("tcp", DB_API_ADDR)
	if err != nil {
		panic(err)
	}

	log.Printf("[db] serving api on %s", DB_API_ADDR)
	server.Serve(listener)
}
