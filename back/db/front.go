package main

import (
	"context"
	"sort"

	"rdbviewer/defs"
	"rdbviewer/shared"
)

func (db *database) Front(context.Context, *shared.FrontReq) (*shared.FrontRes, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := &shared.FrontRes{}

	func() {
		for _, s := range db.data.Shows {
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

		res.LastShows = res.LastShows[0:3]
	}()

	func() {
		for _, b := range db.data.Bootlegs {
			res.LastBootlegs = append(res.LastBootlegs, b)
		}

		sort.Slice(res.LastBootlegs, func(i, j int) bool {
			if shared.PbtimeToTime(res.LastBootlegs[i].FirstSeen) != shared.PbtimeToTime(res.LastBootlegs[j].FirstSeen) {
				return shared.PbtimeToTime(res.LastBootlegs[i].FirstSeen).After(shared.PbtimeToTime(res.LastBootlegs[j].FirstSeen))
			}
			return res.LastBootlegs[i].Id < res.LastBootlegs[j].Id
		})

		res.LastBootlegs = res.LastBootlegs[0:3]

		res.LastBootlegShows = make(map[string]*defs.RadioShow)
		for _, b := range res.LastBootlegs {
			res.LastBootlegShows[b.Show] = db.data.Shows[b.Show]
		}
	}()

	res.Stats = &db.data.Stats

	return res, nil
}
