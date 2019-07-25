package shared

import (
	"fmt"
	"rdbviewer/defs"
	"strings"
)

func LabelTour(s *defs.RadioShow) string {
	switch s.Tour {
	// rh
	case "oaf":
		return "On a Friday"
	case "pbh":
		return "Pablo Honey tour"
	case "bends":
		return "The Bends tour"
	case "okc":
		return "OK Computer tour"
	case "kida":
		return "Kid A tour"
	case "amn":
		return "Amnesiac tour"
	case "httt":
		return "Hail to the Thief tour"
	case "irb":
		return "In Rainbows tour"
	case "tkol":
		return "The King of Kimbs tour"
	case "amsp":
		return "A Moon Shaped Pool tour"

	// jonny / junun
	case "varj":
		return "Various Jonny"
	case "junun":
		return "Junun tour"

	// thom / afp
	case "vart":
		return "Various Thom"
	case "eras":
		return "The Eraser tour"
	case "amok":
		return "Amok tour"
	case "tmb":
		return "Tomorrow's Modern Boxes tour"
	case "anim":
		return "Anima tour"

	// phil
	case "famil":
		return "Familial tour"
	case "weath":
		return "The Weatherhouse tour"
	}
	panic("tour not recognized")
}

func LabelArtist(s *defs.RadioShow) string {
	switch s.Artist {
	case "radiohead":
		return "Radiohead"
	case "thom":
		return "Thom Yorke"
	case "afp":
		return "Atoms for Peace"
	case "jonny":
		return "Jonny Greenwood"
	case "phil":
		return "Phil Selway"
	case "junun":
		return "Junun"
	}
	panic("artist not recognized")
}

func LabelSetlist(url string) string {
	switch url {
	case "citizeninsane":
		return "Citizen Insane"
	case "greenplastic":
		return "Green Plastic Radiohead"
	case "setlistfm":
		return "Setlist.fm"
	case "songkick":
		return "Songkick"
	case "thomthomthom":
		return "Thom thom thom"
	}
	panic("url not recognized")
}

func LabelCountry(s *defs.RadioShow) string {
	return countryMetaData[s.CountryCode].Name
}

func LabelCountryCode(s *defs.RadioShow) string {
	return countryMetaData[s.CountryCode].CodeShortLower
}

func LabelAudioResolution(b *defs.RadioBootleg) string {
	if b.MinfoAudioRate == 0 {
		return "unknown"
	}
	ret := ""
	if b.MinfoAudioDepth != 0 {
		ret += fmt.Sprintf("%dbit", b.MinfoAudioDepth)
	} else {
		ret += "lossy"
	}
	ret += fmt.Sprintf(" %.1fkhz", float64(b.MinfoAudioRate)/1000.0)
	return ret
}

func LabelVideoResolution(b *defs.RadioBootleg) string {
	if b.MinfoVideoHeight == 0 {
		return "unknown"
	}
	return fmt.Sprintf("%dp", b.MinfoVideoHeight)
}

func LabelShortResolution(b *defs.RadioBootleg) string {
	if b.MinfoFormat == "" ||
		(b.Type == "audio" && b.MinfoAudioCodec == "") ||
		(b.Type == "video" && b.MinfoVideoCodec == "") {
		return "unknown"
	}
	switch b.Type {
	case "audio":
		return LabelAudioResolution(b)
	case "video":
		return LabelVideoResolution(b)
	}
	return ""
}

func LabelMediaFormat(b *defs.RadioBootleg) string {
	if b.MinfoFormat == "" {
		return "unknown"
	}
	return b.MinfoFormat
}

func LabelVideoCodec(b *defs.RadioBootleg) string {
	if b.MinfoVideoCodec == "" {
		return "unknown"
	}
	return b.MinfoVideoCodec
}

func LabelAudioCodec(b *defs.RadioBootleg) string {
	if b.MinfoAudioCodec == "" {
		return "unknown"
	}
	return b.MinfoAudioCodec
}

func LabelMediaType(media string) string {
	return strings.Title(media)
}
