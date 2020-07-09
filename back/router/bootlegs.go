package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *router) onDataBootlegs(ctx *gin.Context) {
	var in struct {
		Sort     string
		Text     string
		Media    []string
		AudioRes []string
		VideoRes []string
		CurPage  uint32
	}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&in); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	res, err := h.dbClient.BootlegsFiltered(context.Background(), &shared.BootlegsFilteredReq{
		Sort:     in.Sort,
		Text:     in.Text,
		Media:    in.Media,
		AudioRes: in.AudioRes,
		VideoRes: in.VideoRes,
		CurPage:  in.CurPage,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"fullyLoaded": res.FullyLoaded,
		"choices": func() gin.H {
			if res.Choices != nil {
				return gin.H{
					"media":    res.Choices.Media,
					"audioRes": res.Choices.AudioRes,
					"videoRes": res.Choices.VideoRes,
				}
			}
			return nil
		}(),
		"items": func() (ret []gin.H) {
			for _, b := range res.Items {
				s := res.Shows[b.Show]
				d, _ := time.Parse("2006-01-02", s.Date)

				ret = append(ret, gin.H{
					"id":        b.Id,
					"title":     b.Name,
					"name":      b.Name,
					"firstSeen": formatFirstSeen(b.FirstSeen, "2 Jan 2006"),
					"type":      b.Type,
					"typeLong":  shared.LabelMediaType(b.Type),
					"res":       shared.LabelShortResolution(b),
					"duration": func() string {
						if b.Duration == 0 {
							return ""
						}
						return formatDuration(b.Duration)
					}(),
					"size": humanize.Bytes(b.Size),
					"show": fmt.Sprintf("%s, %s, %s, %s",
						shared.LabelArtist(s), d.Format("2 Jan 2006"), s.City, strings.ToUpper(s.CountryCode)),
					"tour":      s.Tour,
					"labelTour": shared.LabelTour(s),
				})
			}
			return ret
		}(),
	})
}
