package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gin-gonic/gin"

	"rdbviewer/defs"
	"rdbviewer/shared"
)

func magnetLink(f *defs.RadioBootlegFile) string {
	return fmt.Sprintf("magnet:?xt=urn:tree:tiger:%s&xl=%d&dn=%s",
		f.TTH, f.Size, url.QueryEscape(f.Name))
}

func (h *router) onDataBootleg(ctx *gin.Context) {
	res, err := h.dbClient.Bootleg(context.Background(), &shared.BootlegReq{Id: ctx.Param("id")})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	b := res.Item
	s := res.Show
	if b == nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	sd, _ := time.Parse("2006-01-02", s.Date)

	ctx.JSON(http.StatusOK, gin.H{
		"title":            fmt.Sprintf("Bootleg \"%s\"", b.Name),
		"name":             b.Name,
		"type":             b.Type,
		"typeLong":         shared.LabelMediaType(b.Type),
		"size":             humanize.Bytes(b.Size),
		"showArtist":       shared.LabelArtist(s),
		"showDate":         sd.Format("2 January 2006"),
		"showCity":         s.City,
		"showCountry":      strings.ToUpper(s.CountryCode),
		"showId":           s.Id,
		"labelCountryCode": shared.LabelCountryCode(s),
		"labelCountry":     shared.LabelCountry(s),
		"tour":             s.Tour,
		"labelTour":        shared.LabelTour(s),
		"firstSeen":        formatFirstSeen(b.FirstSeen, "2 January 2006"),
		"minfoFormat":      shared.LabelMediaFormat(b),
		"minfoVideoCodec":  shared.LabelVideoCodec(b),
		"minfoVideoRes":    shared.LabelVideoResolution(b),
		"minfoAudioCodec":  shared.LabelAudioCodec(b),
		"minfoAudioRes":    shared.LabelAudioResolution(b),
		"finfo":            b.Finfo,
		"duration": func() string {
			if b.Duration == 0 {
				return "unknown"
			}
			return formatDuration(b.Duration)
		}(),
		"files": func() (ret []gin.H) {
			for _, f := range b.Files {
				ret = append(ret, gin.H{
					"name": f.Name,
					"size": humanize.Bytes(f.Size),
					"duration": func() string {
						if f.Duration == 0 {
							return ""
						}
						return formatDuration(f.Duration)
					}(),
					"TTH":    f.TTH,
					"magnet": template.URL(magnetLink(f)),
				})
			}
			return ret
		}(),
	})
}
