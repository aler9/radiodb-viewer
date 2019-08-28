package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"rdbviewer/shared"
)

func (db *Database) ShowsFiltered(ctx context.Context, in *shared.ShowsFilteredReq) (*shared.ShowsFilteredRes, error) {
	res := &shared.ShowsFilteredRes{
		Choices: &shared.ShowsFilteredChoices{
			Artist:  make(map[string]string),
			Tour:    make(map[string]string),
			Year:    make([]uint32, 2),
			Country: make(map[string]string),
			Media:   make(map[string]string),
		},
	}

	textKeywords := GetTextKeywords(in.Text, 1)

	for _, s := range db.db.Shows {
		// choices
		if _, ok := res.Choices.Artist[s.Artist]; !ok {
			res.Choices.Artist[s.Artist] = shared.LabelArtist(s)
		}
		if _, ok := res.Choices.Tour[s.Tour]; !ok {
			res.Choices.Tour[s.Tour] = shared.LabelTour(s)
		}
		d, _ := time.Parse("2006-01-02", s.Date)
		year := shared.Atoui32(d.Format("2006"))
		if res.Choices.Year[0] == 0 || year < res.Choices.Year[0] {
			res.Choices.Year[0] = year
		}
		if res.Choices.Year[1] == 0 || year > res.Choices.Year[1] {
			res.Choices.Year[1] = year
		}
		if _, ok := res.Choices.Country[s.CountryCode]; !ok {
			res.Choices.Country[s.CountryCode] = shared.LabelCountry(s)
		}
		if s.AudioCount > 0 {
			if _, ok := res.Choices.Media["audio"]; !ok {
				res.Choices.Media["audio"] = shared.LabelMediaType("audio")
			}
		}
		if s.VideoCount > 0 {
			if _, ok := res.Choices.Media["video"]; !ok {
				res.Choices.Media["video"] = shared.LabelMediaType("video")
			}
		}
		if s.MiscCount > 0 {
			if _, ok := res.Choices.Media["misc"]; !ok {
				res.Choices.Media["misc"] = shared.LabelMediaType("misc")
			}
		}

		// filters
		if len(textKeywords) > 0 && func() bool {
			for tw := range textKeywords {
				if func() bool {
					for kw := range db.showKeywords[s.Id] {
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
		if len(in.Artist) > 0 && func() bool {
			for _, v := range in.Artist {
				if v == s.Artist {
					return false
				}
			}
			return true
		}() {
			continue
		}
		if len(in.Media) > 0 && func() bool {
			for _, m := range in.Media {
				switch m {
				case "audio":
					if s.AudioCount > 0 {
						return false
					}
				case "video":
					if s.VideoCount > 0 {
						return false
					}
				case "misc":
					if s.MiscCount > 0 {
						return false
					}
				}
			}
			return true
		}() {
			continue
		}
		if len(in.Tour) > 0 && func() bool {
			for _, v := range in.Tour {
				if s.Tour == v {
					return false
				}
			}
			return true
		}() {
			continue
		}
		if len(in.Year) == 2 && func() bool {
			if year >= in.Year[0] && year <= in.Year[1] {
				return false
			}
			return true
		}() {
			continue
		}
		if len(in.Country) > 0 && func() bool {
			for _, v := range in.Country {
				if s.CountryCode == v {
					return false
				}
			}
			return true
		}() {
			continue
		}

		res.Items = append(res.Items, s)
	}

	switch in.Sort {
	case "date_asc":
		sort.Slice(res.Items, func(i, j int) bool {
			if res.Items[i].Date != res.Items[j].Date {
				return res.Items[i].Date < res.Items[j].Date
			}
			if res.Items[i].Artist != res.Items[j].Artist {
				return res.Items[i].Artist < res.Items[j].Artist
			}
			return res.Items[i].Id < res.Items[j].Id
		})

	default:
		sort.Slice(res.Items, func(i, j int) bool {
			if res.Items[i].Date != res.Items[j].Date {
				return res.Items[i].Date > res.Items[j].Date
			}
			if res.Items[i].Artist != res.Items[j].Artist {
				return res.Items[i].Artist < res.Items[j].Artist
			}
			return res.Items[i].Id < res.Items[j].Id
		})
	}

	start, end, FullyLoaded, ok := Pagination(in.CurPage, uint32(len(res.Items)), 20)
	if !ok {
		return nil, fmt.Errorf("invalid page")
	}
	res.FullyLoaded = FullyLoaded
	res.Items = res.Items[start:end]

	if in.CurPage != 0 {
		res.Choices = nil
	}

	return res, nil
}

func (db *Database) Show(ctx context.Context, req *shared.ShowReq) (*shared.ShowRes, error) {
	res := &shared.ShowRes{}
	if s, ok := db.db.Shows[req.Id]; ok {
		res.Item = s
		for _, b := range db.db.Bootlegs {
			if b.Show == s.Id {
				res.Bootlegs = append(res.Bootlegs, b)
			}
		}
	}
	return res, nil
}
