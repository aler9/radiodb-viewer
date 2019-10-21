package main

import (
	"context"
	"encoding/json"
	"html/template"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *Router) onPageStats(c *gin.Context) {
	stats, err := h.dbClient.Stats(context.Background(), &shared.StatsReq{})
	if err != nil {
		GinServerErrorText(c)
		return
	}

	de, _ := time.Parse("2006-01-02", stats.Stats.DateEarliest)
	dl, _ := time.Parse("2006-01-02", stats.Stats.DateLatest)
	dlb, _ := time.Parse("2006-01-02", stats.Stats.DateLastBootleg)
	g, _ := time.Parse(time.RFC3339, stats.Stats.Generated)

	GinTpl(c, h.frameWrapper(c, FrameConf{
		Title: "Statistics",
		Content: TplRender(h.templates["stats"], gin.H{
			"Stats":           stats.Stats,
			"Generated":       g.Format("2 January 2006"),
			"DateLastBootleg": dlb.Format("2 January 2006"),
			"DateEarliest":    de.Format("2 January 2006"),
			"DateLatest":      dl.Format("2 January 2006"),
			"ShareUniqueSize": humanize.Bytes(stats.Stats.ShareUniqueSize),
			"ShareSize":       humanize.Bytes(stats.Stats.ShareSize),
			"PerYear": func() interface{} {
				ret := []interface{}{
					stats.PerYearShows,
					stats.PerYearBootlegs,
					stats.PerYearSize,
				}
				byt, _ := json.Marshal(ret)
				return template.URL(string(byt))
			}(),
		}),
	}))
}
