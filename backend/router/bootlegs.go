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

func (h *router) onPageBootlegs(c *gin.Context) {
	ginTemplate(c, h.frameWrapper(c, FrameConf{
		Title:   "Bootlegs",
		Content: templateRender(h.templates["bootlegs"], nil),
	}))
}

func (h *router) onDataBootlegs(c *gin.Context) {
	var in struct {
		Sort     string
		Text     string
		Media    []string
		AudioRes []string
		VideoRes []string
		CurPage  uint32
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&in); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
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
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{
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
					"props": gin.H{
						"key":       b.Id,
						"href":      "/bootlegs/" + b.Id,
						"className": "entry",
						"title":     b.Name,
					},
					"cnt": templateRender(h.templates["bootlegentry"], gin.H{
						"Id":        b.Id,
						"Name":      b.Name,
						"FirstSeen": formatFirstSeen(b.FirstSeen, "2 Jan 2006"),
						"Type":      b.Type,
						"TypeLong":  shared.LabelMediaType(b.Type),
						"Res":       shared.LabelShortResolution(b),
						"Duration": func() string {
							if b.Duration == 0 {
								return ""
							}
							return formatDuration(b.Duration)
						}(),
						"Size": humanize.Bytes(b.Size),
						"Show": fmt.Sprintf("%s, %s, %s, %s",
							shared.LabelArtist(s), d.Format("2 Jan 2006"), s.City, strings.ToUpper(s.CountryCode)),
						"Tour":      s.Tour,
						"LabelTour": shared.LabelTour(s),
					}),
				})
			}
			return ret
		}(),
	})
}
