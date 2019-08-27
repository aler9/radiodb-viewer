package main

import (
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"rdbviewer/defs"
	"rdbviewer/shared"
	"time"

	_ "net/http/pprof"
)

const (
	DB_ADDR = "localhost:4002"
	DB_FILE = "/data/radiodb.json"
)

type Database struct {
	db defs.RadioOut
	// derivated, cached data
	showKeywords    map[string]map[string]struct{}
	bootlegKeywords map[string]map[string]struct{}
	perYearShows    map[uint32]int32
	perYearBootlegs map[uint32]int32
	perYearSize     map[uint32]uint64
}

func main() {
	pprofMux := http.DefaultServeMux
	go func() {
		(&http.Server{
			Addr:    ":9998",
			Handler: pprofMux,
		}).ListenAndServe()
	}()
	http.DefaultServeMux = http.NewServeMux()

	rand.Seed(time.Now().UnixNano())

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags)

	db := &Database{}

	MustImportJson(DB_FILE, &db.db)

	db.showKeywords = make(map[string]map[string]struct{})
	for _, s := range db.db.Shows {
		db.showKeywords[s.Id] = make(map[string]struct{})
		pushText := func(text string) {
			for word := range GetTextKeywords(text, 2) {
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
	for _, b := range db.db.Bootlegs {
		db.bootlegKeywords[b.Id] = make(map[string]struct{})
		pushText := func(text string) {
			for word := range GetTextKeywords(text, 2) {
				db.bootlegKeywords[b.Id][word] = struct{}{}
			}
		}

		pushText(b.Name)

		s := db.db.Shows[b.Show]
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
	for _, s := range db.db.Shows {
		d, _ := time.Parse("2006-01-02", s.Date)
		year := shared.Atoui32(d.Format("2006"))
		db.perYearShows[year]++
	}
	for _, b := range db.db.Bootlegs {
		s := db.db.Shows[b.Show]
		d, _ := time.Parse("2006-01-02", s.Date)
		year := shared.Atoui32(d.Format("2006"))
		db.perYearBootlegs[year]++
		db.perYearSize[year] += b.Size
	}

	// init server
	listener, err := net.Listen("tcp", DB_ADDR)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	shared.RegisterDatabaseServer(server, db)
	db.Log("serving database on %s", DB_ADDR)
	server.Serve(listener)
}

func (db *Database) Log(text string, args ...interface{}) {
	log.Printf(text, args...)
}
