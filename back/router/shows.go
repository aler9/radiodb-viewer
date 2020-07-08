package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"rdbviewer/shared"
)

func (h *router) onDataShows(ctx *gin.Context) {
	var in struct {
		Sort    string
		Text    string
		Artist  []string
		Tour    []string
		Year    []uint32
		Country []string
		Media   []string
		CurPage uint32
	}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&in); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	res, err := h.dbClient.ShowsFiltered(context.Background(), &shared.ShowsFilteredReq{
		Sort:    in.Sort,
		Text:    in.Text,
		Artist:  in.Artist,
		Tour:    in.Tour,
		Year:    in.Year,
		Country: in.Country,
		Media:   in.Media,
		CurPage: in.CurPage,
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
					"artist":  res.Choices.Artist,
					"tour":    res.Choices.Tour,
					"year":    res.Choices.Year,
					"country": res.Choices.Country,
					"media":   res.Choices.Media,
				}
			}
			return nil
		}(),
		"items": func() (ret []gin.H) {
			for _, s := range res.Items {
				d, _ := time.Parse("2006-01-02", s.Date)
				ret = append(ret, gin.H{
					"id": s.Id,
					"title": fmt.Sprintf("%s, %s, %s, %s",
						shared.LabelArtist(s), d.Format("2 January 2006"), s.City, shared.LabelCountry(s)),
					"date":             d.Format("2 January 2006"),
					"artist":           s.Artist,
					"city":             s.City,
					"country":          shared.LabelCountry(s),
					"countryCode":      s.CountryCode,
					"labelCountryCode": shared.LabelCountryCode(s),
					"audioCount":       s.AudioCount,
					"videoCount":       s.VideoCount,
					"miscCount":        s.MiscCount,
					"tour":             s.Tour,
					"labelTour":        shared.LabelTour(s),
					"artistLong":       shared.LabelArtist(s),
				})
			}
			return ret
		}(),
	})
}
