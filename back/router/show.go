package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *router) onDataShow(ctx *gin.Context) {
	res, err := h.dbClient.Show(context.Background(), &shared.ShowReq{Id: ctx.Param("id")})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	s := res.Item
	bs := res.Bootlegs
	if s == nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	sort.Slice(bs, func(i, j int) bool {
		return bs[i].Size > bs[j].Size
	})

	d, _ := time.Parse("2006-01-02", s.Date)

	urls := func() (ret []gin.H) {
		for k, v := range s.Url {
			ret = append(ret, gin.H{
				"url":  v,
				"name": shared.LabelSetlist(k),
			})
		}
		return ret
	}()

	sort.Slice(urls, func(i, j int) bool {
		return urls[i]["name"].(string) < urls[j]["name"].(string)
	})

	ctx.JSON(http.StatusOK, gin.H{
		"title": fmt.Sprintf("%s, %s, %s, %s",
			shared.LabelArtist(s), d.Format("2 January 2006"), s.City, shared.LabelCountry(s)),
		"date":             d.Format("2 January 2006"),
		"artistLong":       shared.LabelArtist(s),
		"city":             s.City,
		"countryCode":      s.CountryCode,
		"labelCountryCode": shared.LabelCountryCode(s),
		"labelCountry":     shared.LabelCountry(s),
		"urls":             urls,
		"tour":             s.Tour,
		"labelTour":        shared.LabelTour(s),
		"bootlegs": func() (ret []gin.H) {
			for _, b := range bs {
				ret = append(ret, gin.H{
					"id":       b.Id,
					"name":     b.Name,
					"type":     b.Type,
					"typeLong": shared.LabelMediaType(b.Type),
					"res":      shared.LabelShortResolution(b),
					"duration": func() string {
						if b.Duration == 0 {
							return ""
						}
						return formatDuration(b.Duration)
					}(),
					"size":      humanize.Bytes(b.Size),
					"firstSeen": formatFirstSeen(b.FirstSeen, "2 Jan 2006"),
				})
			}
			return ret
		}(),
	})
}
