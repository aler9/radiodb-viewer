package main

import (
    "sort"
    "context"
    "rdbviewer/back/shared"
    "rdbviewer/back/defs"
)

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
