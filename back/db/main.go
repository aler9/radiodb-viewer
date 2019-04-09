package main

import (
    "os"
    "log"
    "net"
    "time"
    "sort"
    "context"
    "math/rand"
    "google.golang.org/grpc"
    "rdbviewer/back/shared"
    "rdbviewer/back/defs"
)

const(
    DB_ADDR = "localhost:4002"
    DB_FILE = "/data/radiodb.json"
)

type Database struct {
    db              defs.RadioOut
    // derivated, cached data
    showKeywords    map[string]map[string]struct{}
    bootlegKeywords map[string]map[string]struct{}
    perYearShows    map[uint32]int32
    perYearBootlegs map[uint32]int32
    perYearSize     map[uint32]uint64
}

func main() {
    rand.Seed(time.Now().UnixNano())

    log.SetOutput(os.Stdout)
    log.SetFlags(log.LstdFlags)

    db := &Database{}

    MustImportJson(DB_FILE, &db.db)

    db.showKeywords = make(map[string]map[string]struct{})
    for _,s := range db.db.Shows {
        db.showKeywords[s.Id] = make(map[string]struct{})
        pushText := func(text string) {
            for word,_ := range shared.GetTextKeywords(text, 2) {
                db.showKeywords[s.Id][word] = struct{}{}
            }
        }
        pushText(shared.ArtistLabel(s))
        pushText(s.City)
        pushText(s.CountryCode)
        pushText(shared.CountryLabel(s))
        d,_ := time.Parse("2006-01-02", s.Date)
        pushText(d.Format("01 02 2006 January"))
    }

    db.bootlegKeywords = make(map[string]map[string]struct{})
    for _,b := range db.db.Bootlegs {
        db.bootlegKeywords[b.Id] = make(map[string]struct{})
        pushText := func(text string) {
            for word,_ := range shared.GetTextKeywords(text, 2) {
                db.bootlegKeywords[b.Id][word] = struct{}{}
            }
        }

        pushText(b.Name)

        s := db.db.Shows[b.Show]
        pushText(shared.ArtistLabel(s))
        pushText(s.City)
        pushText(s.CountryCode)
        pushText(shared.CountryLabel(s))
        d,_ := time.Parse("2006-01-02", s.Date)
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
    for _,s := range db.db.Shows {
        d,_ := time.Parse("2006-01-02", s.Date)
        year := shared.Atoui32(d.Format("2006"))
        db.perYearShows[year]++
    }
    for _,b := range db.db.Bootlegs {
        s := db.db.Shows[b.Show]
        d,_ := time.Parse("2006-01-02", s.Date)
        year := shared.Atoui32(d.Format("2006"))
        db.perYearBootlegs[year]++
        db.perYearSize[year] += b.Size
    }

    // init server
    listener,err := net.Listen("tcp", DB_ADDR)
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

func (db *Database) Front(context.Context, *shared.FrontReq) (*shared.FrontRes,error) {
    res := &shared.FrontRes{}

    func() {
        for _,s := range db.db.Shows {
            res.LastShows = append(res.LastShows, s)
        }

        sort.Slice(res.LastShows, func(i, j int) bool {
            if res.LastShows[i].Date != res.LastShows[j].Date {
                return res.LastShows[i].Date > res.LastShows[j].Date
            }
            if res.LastShows[i].Artist != res.LastShows[j].Artist {
                return res.LastShows[i].Artist < res.LastShows[j].Artist
            }
            return res.LastShows[i].Id < res.LastShows[j].Id
        })

        res.LastShows = res.LastShows[ 0 : 3 ]
    }()

    func() {
        for _,b := range db.db.Bootlegs {
            res.LastBootlegs = append(res.LastBootlegs, b)
        }

        sort.Slice(res.LastBootlegs, func(i, j int) bool {
            if shared.PbtimeToTime(res.LastBootlegs[i].FirstSeen) != shared.PbtimeToTime(res.LastBootlegs[j].FirstSeen) {
                return shared.PbtimeToTime(res.LastBootlegs[i].FirstSeen).After(shared.PbtimeToTime(res.LastBootlegs[j].FirstSeen))
            }
            return res.LastBootlegs[i].Id < res.LastBootlegs[j].Id
        })

        res.LastBootlegs = res.LastBootlegs[ 0 : 3 ]

        res.LastBootlegShows = make(map[string]*defs.RadioShow)
        for _,b := range res.LastBootlegs {
            res.LastBootlegShows[b.Show] = db.db.Shows[b.Show]
        }
    }()

    res.Stats = &db.db.Stats

    return res, nil
}

func (db *Database) Stats(context.Context, *shared.StatsReq) (*shared.StatsRes,error) {
    res := &shared.StatsRes{}
    res.Stats = &db.db.Stats
    res.PerYearShows = db.perYearShows
    res.PerYearBootlegs = db.perYearBootlegs
    res.PerYearSize = db.perYearSize
    return res, nil
}
