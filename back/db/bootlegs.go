package main

import (
    "fmt"
    "sort"
    "strings"
    "context"
    "math/rand"
    "rdbviewer/back/shared"
    "rdbviewer/back/defs"
)

func (db *Database) BootlegsFiltered(ctx context.Context, in *shared.BootlegsFilteredReq) (*shared.BootlegsFilteredRes,error) {
    res := &shared.BootlegsFilteredRes{
        Choices: &shared.BootlegsFilteredChoices{
            Media: make(map[string]string),
            AudioRes: make(map[string]string),
            VideoRes: make(map[string]string),
        },
        Shows: make(map[string]*defs.RadioShow),
    }

    textKeywords := GetTextKeywords(in.Text, 1)

    for _,b := range db.db.Bootlegs {
        // choices
        if _,ok := res.Choices.Media[b.Type]; !ok {
            res.Choices.Media[b.Type] = shared.MediaTypeLabel(b.Type)
        }
        var audioRes string
        if b.Type == "audio" {
            audioRes = shared.AudioResolution(b)
            if _,ok := res.Choices.AudioRes[audioRes]; !ok {
                res.Choices.AudioRes[audioRes] = audioRes
            }
        }
        var videoRes string
        if b.Type == "video" {
            videoRes = shared.VideoResolution(b) //fmt.Sprintf("%dp", b.MinfoVideoHeight)
            if _,ok := res.Choices.VideoRes[videoRes]; !ok {
                res.Choices.VideoRes[videoRes] = videoRes
            }
        }

        // filters
        if len(textKeywords) > 0 && func() bool {
            // search by TTH
            if len(textKeywords) == 1 && len(FirstKey(textKeywords)) == 39 {
                tth := strings.ToUpper(FirstKey(textKeywords))

                for _,f := range b.Files {
                    if f.TTH == tth {
                        return false
                    }
                }
                return true
            }

            // normal search
            for tw,_ := range textKeywords {
                if func() bool {
                    for kw,_ := range db.bootlegKeywords[b.Id] {
                        if strings.Contains(kw, tw) {
                            return false
                        }
                    }
                    return true
                }() {
                    return true
                }
            }
            return false
        }() {
            continue
        }
        if len(in.Media) > 0 && func() bool {
            for _,r := range in.Media {
                if r == b.Type {
                    return false
                }
            }
            return true
        }() {
            continue
        }
        if len(in.AudioRes) > 0 && (b.Type != "audio" || func() bool {
            for _,r := range in.AudioRes {
                if r == audioRes {
                    return false
                }
            }
            return true
        }()) {
            continue
        }
        if len(in.VideoRes) > 0 && (b.Type != "video" || func() bool {
            for _,r := range in.VideoRes {
                if r == videoRes {
                    return false
                }
            }
            return true
        }()) {
            continue
        }

        res.Items = append(res.Items, b)
    }

    switch in.Sort {
    case "sdate_desc":
        sort.Slice(res.Items, func(i, j int) bool {
            if d1, d2 := db.db.Shows[res.Items[i].Show].Date, db.db.Shows[res.Items[j].Show].Date; d1 != d2 {
                return d1 > d2
            }
            return res.Items[i].Id < res.Items[j].Id
        })

    case "sdate_asc":
        sort.Slice(res.Items, func(i, j int) bool {
            if d1, d2 := db.db.Shows[res.Items[i].Show].Date, db.db.Shows[res.Items[j].Show].Date; d1 != d2 {
                return d1 < d2
            }
            return res.Items[i].Id < res.Items[j].Id
        })

    case "size_desc":
        sort.Slice(res.Items, func(i, j int) bool {
            if res.Items[i].Size != res.Items[j].Size {
                return res.Items[i].Size > res.Items[j].Size
            }
            return res.Items[i].Id < res.Items[j].Id
        })

    case "size_asc":
        sort.Slice(res.Items, func(i, j int) bool {
            if res.Items[i].Size != res.Items[j].Size {
                return res.Items[i].Size < res.Items[j].Size
            }
            return res.Items[i].Id < res.Items[j].Id
        })

    default:
        sort.Slice(res.Items, func(i, j int) bool {
            if shared.PbtimeToTime(res.Items[i].FirstSeen) != shared.PbtimeToTime(res.Items[j].FirstSeen) {
                return shared.PbtimeToTime(res.Items[i].FirstSeen).After(shared.PbtimeToTime(res.Items[j].FirstSeen))
            }
            return res.Items[i].Id < res.Items[j].Id
        })
    }

    start,end,FullyLoaded,ok := Pagination(in.CurPage, uint32(len(res.Items)), 20)
    if !ok {
        return nil, fmt.Errorf("invalid page")
    }
    res.FullyLoaded = FullyLoaded
    res.Items = res.Items[ start : end ]

    for _,b := range res.Items {
        res.Shows[b.Show] = db.db.Shows[b.Show]
    }

    if in.CurPage != 0 {
        res.Choices = nil
    }

    return res, nil
}

func (db *Database) Bootleg(ctx context.Context, req *shared.BootlegReq) (*shared.BootlegRes,error) {
    res := &shared.BootlegRes{}
    if b,ok := db.db.Bootlegs[req.Id]; ok {
        if s,ok := db.db.Shows[b.Show]; ok {
            res.Item = b
            res.Show = s
        }
    }
    return res, nil
}

func (db *Database) BootlegRand(context.Context, *shared.BootlegRandReq) (*shared.BootlegRandRes,error) {
    res := &shared.BootlegRandRes{}

    var bs []string
    for _,b := range db.db.Bootlegs {
        bs = append(bs, b.Id)
    }
    res.Id = bs[rand.Intn(len(bs))]

    return res, nil
}
