package main

import (
    "fmt"
    "sort"
    "time"
    "context"
    "github.com/gin-gonic/gin"
    "github.com/dustin/go-humanize"
    "rdbviewer/back/shared"
)

func (h *Router) onPageShow(c *gin.Context) {
    res,err := h.dbClient.Show(context.Background(), &shared.ShowReq{ Id: c.Param("id") })
    if err != nil {
        GinServerErrorText(c)
        return
    }

    s := res.Item
    bs := res.Bootlegs
    if s == nil {
        GinNotFoundText(c)
        return
    }

    sort.Slice(bs, func(i, j int) bool {
        return bs[i].Size > bs[j].Size
    })

    d,_ := time.Parse("2006-01-02", s.Date)

    urls := func() (ret []gin.H) {
        for k,v := range s.Url {
            ret = append(ret, gin.H{
                "Url": v,
                "Name": shared.SetlistUrlLabel(k),
            })
        }
        return ret
    }()

    sort.Slice(urls, func(i, j int) bool {
        return urls[i]["Name"].(string) < urls[j]["Name"].(string)
    })

    GinTpl(c, h.frameWrapper(c, FrameConf{
        Title: fmt.Sprintf("%s, %s, %s, %s",
            shared.ArtistLabel(s), d.Format("2 January 2006"), s.City, shared.CountryLabel(s)),
        Content: TplExecute(h.templates["show"], gin.H{
            "Date": d.Format("2 January 2006"),
            "ArtistLong": shared.ArtistLabel(s),
            "City": s.City,
            "CountryCode": s.CountryCode,
            "CountryCodeShort": shared.CountryCodeShort(s),
            "CountryLabel": shared.CountryLabel(s),
            "Urls": urls,
            "Tour": s.Tour,
            "TourLabel": shared.TourLabel(s),
            "Bootlegs": func() (ret []gin.H) {
                for _,b := range bs {
                    ret = append(ret, gin.H{
                        "Id": b.Id,
                        "Name": b.Name,
                        "Type": b.Type,
                        "TypeLong": shared.MediaTypeLabel(b.Type),
                        "Res": shared.ShortResolution(b),
                        "Duration": func() string {
                            if b.Duration == 0 {
                                return ""
                            }
                            return FormatDuration(b.Duration)
                        }(),
                        "Size": humanize.Bytes(b.Size),
                        "FirstSeen": FormatFirstSeen(b.FirstSeen, "2 Jan 2006"),
                    })
                }
                return ret
            }(),
        }),
    }))
}
