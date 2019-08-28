package main

import (
	"context"

	"rdbviewer/shared"
)

func (db *Database) Stats(context.Context, *shared.StatsReq) (*shared.StatsRes, error) {
	res := &shared.StatsRes{}
	res.Stats = &db.db.Stats
	res.PerYearShows = db.perYearShows
	res.PerYearBootlegs = db.perYearBootlegs
	res.PerYearSize = db.perYearSize
	return res, nil
}
