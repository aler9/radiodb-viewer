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

func (h *router) onPageShows(c *gin.Context) {
	ginTemplate(c, h.frameWrapper(c, FrameConf{
		Title:   "Shows",
		Content: templateRender(h.templates["shows"], nil),
	}))
}

func (h *router) onDataShows(c *gin.Context) {
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
	if err := json.NewDecoder(c.Request.Body).Decode(&in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
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
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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
					"props": gin.H{
						"key":       s.Id,
						"href":      "/shows/" + s.Id,
						"className": "entry",
						"title": fmt.Sprintf("%s, %s, %s, %s",
							shared.LabelArtist(s), d.Format("2 January 2006"), s.City, shared.LabelCountry(s)),
					},
					"cnt": templateRender(h.templates["showentry"], gin.H{
						"Id":               s.Id,
						"Date":             d.Format("2 January 2006"),
						"Artist":           s.Artist,
						"City":             s.City,
						"Country":          shared.LabelCountry(s),
						"CountryCode":      s.CountryCode,
						"LabelCountryCode": shared.LabelCountryCode(s),
						"AudioCount":       s.AudioCount,
						"VideoCount":       s.VideoCount,
						"MiscCount":        s.MiscCount,
						"Tour":             s.Tour,
						"LabelTour":        shared.LabelTour(s),
						"ArtistLong":       shared.LabelArtist(s),
					}),
				})
			}
			return ret
		}(),
	})
}