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

func MagnetLink(f *defs.RadioBootlegFile) string {
	return fmt.Sprintf("magnet:?xt=urn:tree:tiger:%s&xl=%d&dn=%s",
		f.TTH,
		f.Size,
		url.QueryEscape(f.Name))
}

func SafeHtmlNewlines(in string) template.HTML {
	ret := template.HTMLEscapeString(in)
	ret = strings.Replace(ret, "\n", "<br />", -1)
	return template.HTML(ret)
}

func (h *router) onPageBootleg(c *gin.Context) {
	res, err := h.dbClient.Bootleg(context.Background(), &shared.BootlegReq{Id: c.Param("id")})
	if err != nil {
		http.Error(c.Writer, "500 internal server error", http.StatusInternalServerError)
		return
	}

	b := res.Item
	s := res.Show
	if b == nil {
		http.Error(c.Writer, "404 page not found", http.StatusNotFound)
		return
	}

	sd, _ := time.Parse("2006-01-02", s.Date)

	ginTemplate(c, h.frameWrapper(c, FrameConf{
		Title: fmt.Sprintf("Bootleg \"%s\"", b.Name),
		Content: templateRender(h.templates["bootleg"], gin.H{
			"Name":             b.Name,
			"Type":             b.Type,
			"TypeLong":         shared.LabelMediaType(b.Type),
			"Size":             humanize.Bytes(b.Size),
			"ShowArtist":       shared.LabelArtist(s),
			"ShowDate":         sd.Format("2 January 2006"),
			"ShowCity":         s.City,
			"ShowCountry":      strings.ToUpper(s.CountryCode),
			"ShowUrl":          "/shows/" + s.Id,
			"LabelCountryCode": shared.LabelCountryCode(s),
			"LabelCountry":     shared.LabelCountry(s),
			"Tour":             s.Tour,
			"LabelTour":        shared.LabelTour(s),
			"FirstSeen":        formatFirstSeen(b.FirstSeen, "2 January 2006"),
			"MinfoFormat":      shared.LabelMediaFormat(b),
			"MinfoVideoCodec":  shared.LabelVideoCodec(b),
			"MinfoVideoRes":    shared.LabelVideoResolution(b),
			"MinfoAudioCodec":  shared.LabelAudioCodec(b),
			"MinfoAudioRes":    shared.LabelAudioResolution(b),
			"Finfo":            SafeHtmlNewlines(b.Finfo),
			"Duration": func() string {
				if b.Duration == 0 {
					return "unknown"
				}
				return formatDuration(b.Duration)
			}(),
			"Files": func() (ret []gin.H) {
				for _, f := range b.Files {
					ret = append(ret, gin.H{
						"Name": f.Name,
						"Size": humanize.Bytes(f.Size),
						"Duration": func() string {
							if f.Duration == 0 {
								return ""
							}
							return formatDuration(f.Duration)
						}(),
						"TTH":    f.TTH,
						"Magnet": template.URL(MagnetLink(f)),
					})
				}
				return ret
			}(),
		}),
	}))
}