package main

import (
	"context"

	"rdbviewer/shared"
)

func (db *database) Stats(context.Context, *shared.StatsReq) (*shared.StatsRes, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := &shared.StatsRes{}
	res.Stats = &db.data.Stats
	res.PerYearShows = db.perYearShows
	res.PerYearBootlegs = db.perYearBootlegs
	res.PerYearSize = db.perYearSize
	return res, nil
}
