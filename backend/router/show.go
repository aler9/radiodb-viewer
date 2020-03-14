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

func (h *router) onPageShow(c *gin.Context) {
	res, err := h.dbClient.Show(context.Background(), &shared.ShowReq{Id: c.Param("id")})
	if err != nil {
		http.Error(c.Writer, "500 internal server error", http.StatusInternalServerError)
		return
	}

	s := res.Item
	bs := res.Bootlegs
	if s == nil {
		http.Error(c.Writer, "404 page not found", http.StatusNotFound)
		return
	}

	sort.Slice(bs, func(i, j int) bool {
		return bs[i].Size > bs[j].Size
	})

	d, _ := time.Parse("2006-01-02", s.Date)

	urls := func() (ret []gin.H) {
		for k, v := range s.Url {
			ret = append(ret, gin.H{
				"Url":  v,
				"Name": shared.LabelSetlist(k),
			})
		}
		return ret
	}()

	sort.Slice(urls, func(i, j int) bool {
		return urls[i]["Name"].(string) < urls[j]["Name"].(string)
	})

	ginTemplate(c, h.frameWrapper(c, FrameConf{
		Title: fmt.Sprintf("%s, %s, %s, %s",
			shared.LabelArtist(s), d.Format("2 January 2006"), s.City, shared.LabelCountry(s)),
		Content: templateRender(h.templates["show"], gin.H{
			"Date":             d.Format("2 January 2006"),
			"ArtistLong":       shared.LabelArtist(s),
			"City":             s.City,
			"CountryCode":      s.CountryCode,
			"LabelCountryCode": shared.LabelCountryCode(s),
			"LabelCountry":     shared.LabelCountry(s),
			"Urls":             urls,
			"Tour":             s.Tour,
			"LabelTour":        shared.LabelTour(s),
			"Bootlegs": func() (ret []gin.H) {
				for _, b := range bs {
					ret = append(ret, gin.H{
						"Id":       b.Id,
						"Name":     b.Name,
						"Type":     b.Type,
						"TypeLong": shared.LabelMediaType(b.Type),
						"Res":      shared.LabelShortResolution(b),
						"Duration": func() string {
							if b.Duration == 0 {
								return ""
							}
							return formatDuration(b.Duration)
						}(),
						"Size":      humanize.Bytes(b.Size),
						"FirstSeen": formatFirstSeen(b.FirstSeen, "2 Jan 2006"),
					})
				}
				return ret
			}(),
		}),
	}))
}
