package main

import (
    "fmt"
    "time"
    "strings"
    "context"
    "net/url"
    "html/template"
    "github.com/gin-gonic/gin"
    "github.com/dustin/go-humanize"
    "rdbviewer/back/shared"
    "rdbviewer/back/defs"
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

func (h *Router) onPageBootleg(c *gin.Context) {
    res,err := h.dbClient.Bootleg(context.Background(), &shared.BootlegReq{ Id: c.Param("id") })
    if err != nil {
        GinServerErrorText(c)
        return
    }

    b := res.Item
    s := res.Show
    if b == nil {
        GinNotFoundText(c)
        return
    }

    sd,_ := time.Parse("2006-01-02", s.Date)

    GinTpl(c, h.frameWrapper(c, FrameConf{
        Title: fmt.Sprintf("Bootleg \"%s\"", b.Name),
        Content: TplExecute(h.templates["bootleg"], gin.H{
            "Name": b.Name,
            "Type": b.Type,
            "TypeLong": shared.MediaTypeLabel(b.Type),
            "Size": humanize.Bytes(b.Size),
            "ShowArtist": shared.ArtistLabel(s),
            "ShowDate": sd.Format("2 January 2006"),
            "ShowCity": s.City,
            "ShowCountry": strings.ToUpper(s.CountryCode),
            "ShowUrl": "/shows/"+s.Id,
            "CountryCodeShort": shared.CountryCodeShort(s),
            "CountryLabel": shared.CountryLabel(s),
            "Tour": s.Tour,
            "TourLabel": shared.TourLabel(s),
            "FirstSeen": shared.FormatFirstSeen(b.FirstSeen, "2 January 2006"),
            "MinfoFormat": shared.FormatLabel(b),
            "MinfoVideoCodec": shared.VideoCodecLabel(b),
            "MinfoVideoRes": shared.VideoResolution(b),
            "MinfoAudioCodec": shared.AudioCodecLabel(b),
            "MinfoAudioRes": shared.AudioResolution(b),
            "Finfo": SafeHtmlNewlines(b.Finfo),
            "Duration": func() string {
                if b.Duration == 0 {
                    return "unknown"
                }
                return shared.FormatDuration(b.Duration)
            }(),
            "Files": func() (ret []gin.H) {
                for _,f := range b.Files {
                    ret = append(ret, gin.H{
                        "Name": f.Name,
                        "Size": humanize.Bytes(f.Size),
                        "Duration": func() string {
                            if f.Duration == 0 {
                                return ""
                            }
                            return shared.FormatDuration(f.Duration)
                        }(),
                        "TTH": f.TTH,
                        "Magnet": template.URL(MagnetLink(f)),
                    })
                }
                return ret
            }(),
        }),
    }))
}
