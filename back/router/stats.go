package main

import (
	"context"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *router) onDataStats(ctx *gin.Context) {
	stats, err := h.dbClient.Stats(context.Background(), &shared.StatsReq{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	de, _ := time.Parse("2006-01-02", stats.Stats.DateEarliest)
	dl, _ := time.Parse("2006-01-02", stats.Stats.DateLatest)
	dlb, _ := time.Parse("2006-01-02", stats.Stats.DateLastBootleg)
	g, _ := time.Parse(time.RFC3339, stats.Stats.Generated)

	ctx.JSON(http.StatusOK, gin.H{
		"stats":           stats.Stats,
		"generated":       g.Format("2 January 2006"),
		"dateLastBootleg": dlb.Format("2 January 2006"),
		"dateEarliest":    de.Format("2 January 2006"),
		"dateLatest":      dl.Format("2 January 2006"),
		"shareUniqueSize": humanize.Bytes(stats.Stats.ShareUniqueSize),
		"shareSize":       humanize.Bytes(stats.Stats.ShareSize),
		"perYear": []interface{}{
			stats.PerYearShows,
			stats.PerYearBootlegs,
			stats.PerYearSize,
		},
	})
}
