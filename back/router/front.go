package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *router) onDataFront(ctx *gin.Context) {
	front, err := h.dbClient.Front(context.Background(), &shared.FrontReq{})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	g, _ := time.Parse(time.RFC3339, front.Stats.Generated)

	limitSize := func(in string) string {
		if len(in) > 55 {
			in = in[:55] + "..."
		}
		return in
	}

	ctx.JSON(http.StatusOK, gin.H{
		"showsLast": func() (ret []gin.H) {
			for _, s := range front.LastShows {
				da, _ := time.Parse("2006-01-02", s.Date)
				ret = append(ret, gin.H{
					"id":      s.Id,
					"title":   fmt.Sprintf("%s, %s", shared.LabelArtist(s), da.Format("2 January 2006")),
					"country": s.City + ", " + strings.ToUpper(s.CountryCode),
					"tour":    s.Tour,
				})
			}
			return ret
		}(),
		"bootlegsLast": func() (ret []gin.H) {
			for _, b := range front.LastBootlegs {
				s := front.LastBootlegShows[b.Show]
				ret = append(ret, gin.H{
					"id":   b.Id,
					"name": limitSize(b.Name),
					"tour": s.Tour,
				})
			}
			return ret
		}(),
		"generated":    g.Format("2 January 2006"),
		"bootlegCount": front.Stats.BootlegCount,
		"showCount":    front.Stats.ShowCount,
		"shareSize":    humanize.Bytes(front.Stats.ShareSize),
	})
}
